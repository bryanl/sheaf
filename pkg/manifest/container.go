/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/util/jsonpath"

	"github.com/bryanl/sheaf/pkg/sheaf"
	"github.com/pivotal/image-relocation/pkg/images"
)

// ContainerImages returns images from containers in manifest path
func ContainerImages(manifestPath string, definedImages []sheaf.UserDefinedImage) (images.Set, error) {
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return images.Empty, fmt.Errorf("read file: %w", err)
	}

	set := images.Empty

	docs, err := manifestDocuments(data)
	if err != nil {
		return images.Empty, fmt.Errorf("read documents from %s: %w", manifestPath, err)
	}

	for _, doc := range docs {
		results, err := jsonPathMultiSearch(doc, "{range ..spec.containers[*]}{.image}{','}{end}")
		if err != nil {
			return images.Empty, fmt.Errorf("json path search: %w", err)
		}

		bufImages, err := images.New(results...)
		if err != nil {
			return images.Empty, err
		}
		set = set.Union(bufImages)

		for _, udi := range definedImages {
			if !(doc["apiVersion"] == udi.APIVersion && doc["kind"] == udi.Kind) {
				continue
			}

			var bufImages images.Set
			switch udi.Type {
			case sheaf.SingleResult:
				result, err := jsonPathSearch(doc, udi.JSONPath)
				if err != nil {
					return images.Empty, fmt.Errorf("user defined image search %q: %w", udi.JSONPath, err)
				}

				bufImages, err = images.New(result)
				if err != nil {
					return images.Empty, err
				}
			case sheaf.MultiResult:
				results, err := jsonPathMultiSearch(doc, udi.JSONPath)
				if err != nil {
					return images.Empty, fmt.Errorf("user defined image search %q: %w", udi.JSONPath, err)
				}

				if results == nil {
					results = []string{}
				}

				bufImages, err = images.New(results...)
				if err != nil {
					return images.Empty, err
				}
			default:
				return images.Empty, fmt.Errorf("user defined image type %q is invalid", udi.Type)
			}

			set = set.Union(bufImages)

		}

	}

	return set, nil
}

func filterEmpty(ss []string) []string {
	var result []string
	for _, s := range ss {
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}

func jsonPathSearch(doc map[string]interface{}, query string) (string, error) {
	j := jsonpath.New("parser")
	if err := j.Parse(query); err != nil {
		return "", fmt.Errorf("unable to parse: %w", err)
	}

	var buf bytes.Buffer
	if err := j.Execute(&buf, doc); err != nil {
		// jsonpath doesn't return a helpful error, so look at the error message string.
		if strings.Contains(err.Error(), "is not found") {
			return "", nil
		}
		return "", fmt.Errorf("search manifest for containers: %w", err)
	}

	return buf.String(), nil
}

func jsonPathMultiSearch(doc map[string]interface{}, query string) ([]string, error) {
	result, err := jsonPathSearch(doc, query)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, nil
	}

	return filterEmpty(strings.Split(result, ",")), nil
}

func manifestDocuments(in []byte) ([]map[string]interface{}, error) {
	r := bytes.NewReader(in)
	decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)

	var list []map[string]interface{}

	for {
		var m map[string]interface{}
		if err := decoder.Decode(&m); err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("decode failed: %w", err)
		}

		list = append(list, m)
	}

	return list, nil
}
