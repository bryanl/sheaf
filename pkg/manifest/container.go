/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/bryanl/sheaf/internal/yamlutil"
	"github.com/bryanl/sheaf/pkg/sheaf"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/images"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
)

// ContainerImagesFromBytes returns container images referenced in manifest bytes.
func ContainerImagesFromBytes(data []byte, userDefinedImages []sheaf.UserDefinedImage) (images.Set, error) {
	set := images.Empty

	docs, err := manifestDocuments(data)
	if err != nil {
		return images.Empty, fmt.Errorf("read documents: %w", err)
	}

	for _, doc := range docs {
		results, err := jsonPathSearch(doc, "..spec.containers[*].image")
		if err != nil {
			return images.Empty, fmt.Errorf("json path search: %w", err)
		}

		bufImages, err := images.New(results...)
		if err != nil {
			return images.Empty, err
		}
		set = set.Union(bufImages)

		for _, udi := range userDefinedImages {
			var d map[string]interface{}
			if err := doc.Decode(&d); err != nil {
				return images.Empty, fmt.Errorf("YAML decode: %w", err)
			}
			if !(d["apiVersion"] == udi.APIVersion && d["kind"] == udi.Kind) {
				continue
			}

			result, err := jsonPathSearch(doc, udi.JSONPath)
			if err != nil {
				return images.Empty, fmt.Errorf("user defined image search %q: %w", udi.JSONPath, err)
			}

			bufImages, err := images.New(result...)
			if err != nil {
				return images.Empty, err
			}

			set = set.Union(bufImages)

		}

	}

	return set, nil
}

// ContainerImages returns images from containers in manifest path
func ContainerImages(manifestPath string, definedImages []sheaf.UserDefinedImage) (images.Set, error) {
	data, err := ioutil.ReadFile(manifestPath)
	if err != nil {
		return images.Empty, fmt.Errorf("read file: %w", err)
	}

	imagesSet, err := ContainerImagesFromBytes(data, definedImages)
	if err != nil {
		return images.Empty, fmt.Errorf("find container images in %s: %w", manifestPath, err)
	}

	return imagesSet, nil
}

// MapContainer applies the mapping to the images in the input manifest and returns the modified manifest
func MapContainer(manifest []byte, userDefinedImages []sheaf.UserDefinedImage, mapping func(originalImage image.Name) (image.Name, error)) ([]byte, error) {
	refMapping := func(originalImage string) (string, error) {
		i, err := image.NewName(originalImage)
		if err != nil {
			return "", err
		}
		mappedImage, err := mapping(i)
		if err != nil {
			return "", err
		}
		return mappedImage.String(), nil
	}

	docs, err := manifestDocuments(manifest)
	if err != nil {
		return nil, fmt.Errorf("read documents: %w", err)
	}

	newDocs := []string{}
	for _, doc := range docs {
		// Skip empty documents
		if doc.Content == nil {
			continue
		}
		imageNodes, err := jsonPathSearchNodes(doc, "..spec.containers[*].image")
		if err != nil {
			return nil, fmt.Errorf("json path search: %w", err)
		}

		for _, i := range imageNodes {
			newValue, err := refMapping(i.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to map image: %w", err)
			}
			i.Value = newValue
		}

		for _, udi := range userDefinedImages {
			var d map[string]interface{}
			if err := doc.Decode(&d); err != nil {
				return nil, fmt.Errorf("YAML decode: %w", err)
			}
			if !(d["apiVersion"] == udi.APIVersion && d["kind"] == udi.Kind) {
				continue
			}

			imageNodes, err := jsonPathSearchNodes(doc, udi.JSONPath)
			if err != nil {
				return nil, fmt.Errorf("user defined image search %q: %w", udi.JSONPath, err)
			}

			for _, i := range imageNodes {
				newValue, err := refMapping(i.Value)
				if err != nil {
					return nil, fmt.Errorf("failed to map image: %w", err)
				}
				i.Value = newValue
			}
		}

		var buf bytes.Buffer
		e := yaml.NewEncoder(&buf)
		e.SetIndent(2)

		err = e.Encode(doc)
		if err != nil {
			return nil, fmt.Errorf("cannot marshal node %#v: %w", doc, err)
		}
		e.Close()

		newDocs = append(newDocs, buf.String())
	}

	return []byte(strings.Join(newDocs, "\n---\n")), nil // add leading newline since splitting can drop newlines
}

func jsonPathSearch(doc *yaml.Node, query string) ([]string, error) {
	imageNodes, err := jsonPathSearchNodes(doc, query)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, imageNode := range imageNodes {
		result = append(result, imageNode.Value)
	}

	return result, nil
}

func jsonPathSearchNodes(doc *yaml.Node, query string) ([]*yaml.Node, error) {
	p, err := yamlpath.NewPath(query)
	if err != nil {
		return nil, fmt.Errorf("unable to parse query: %w", err)
	}

	return p.Find(doc)
}

func manifestDocuments(in []byte) ([]*yaml.Node, error) {
	docs, err := yamlutil.Split(in)
	if err != nil {
		return nil, err
	}

	nodes := []*yaml.Node{}
	for _, doc := range docs {
		var n yaml.Node
		err := yaml.Unmarshal(doc, &n)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, &n)
	}

	return nodes, nil
}
