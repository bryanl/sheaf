package sheaf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/util/jsonpath"
)

// ContainerImages returns containers in manifest path
func ContainerImages(manifestPath string) ([]string, error) {
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	r := bytes.NewReader(data)
	decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)

	containerMap := make(map[string]bool)

	for {
		var m map[string]interface{}
		if err := decoder.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decode failed: %w", err)
		}

		j := jsonpath.New("parser")
		if err := j.Parse("{range ..spec.containers[*]}{.image}{','}{end}"); err != nil {
			return nil, fmt.Errorf("unable to parse: %w", err)
		}

		var buf bytes.Buffer
		if err := j.Execute(&buf, m); err != nil {
			// jsonpath doesn't return a helpful error, so looking at the error message
			if strings.Contains(err.Error(), "is not found") {
				continue
			}
			return nil, fmt.Errorf("search manifest for containers: %w", err)
		}

		for _, s := range strings.Split(buf.String(), ",") {
			if s != "" {
				containerMap[s] = true
			}
		}
	}

	var list []string
	for k := range containerMap {
		list = append(list, k)
	}

	return list, nil
}
