package sheaf

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"k8s.io/client-go/util/jsonpath"
	"sigs.k8s.io/yaml"
)

// Containers returns containers in manifest path
func Containers(manifestPath string) ([]string, error) {
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", manifestPath, err)
	}

	var m map[string]interface{}

	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("decode manifest %q: %w", manifestPath, err)
	}

	j := jsonpath.New("parser")
	if err := j.Parse("{range ..spec.containers[*]}{.image}{','}{end}"); err != nil {
		return nil, fmt.Errorf("unable to parse JSON in %q: %w", manifestPath, err)
	}

	var buf bytes.Buffer
	if err := j.Execute(&buf, m); err != nil {
		// jsonpath doesn't return a helpful error, so looking at the error message
		if strings.Contains(err.Error(), "is not found") {
			return nil, nil
		}
		return nil, fmt.Errorf("search manifest for containers: %w", err)
	}

	containers := strings.Split(buf.String(), ",")
	return compactStringSlice(containers), nil
}

// compactStringSlice removes empty elements from a string slice.
func compactStringSlice(sl []string) []string {
	var out []string
	for i := range sl {
		if sl[i] != "" {
			out = append(out, sl[i])
		}
	}

	return out
}
