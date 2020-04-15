/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"sort"

	"go.uber.org/multierr"
	"k8s.io/client-go/util/jsonpath"
)

// UserDefinedImageType is the type of user defined image.
type UserDefinedImageType string

const (
	// MultiResult states the JSON path will return multiple values.
	MultiResult UserDefinedImageType = "multiple"
	// SingleResult states the JSON path will return a single value.
	SingleResult UserDefinedImageType = "single"
)

// UserDefinedImageTypes is a list of user defined image types as a string.
var UserDefinedImageTypes = []string{string(MultiResult), string(SingleResult)}

// UserDefinedImage is a user defined image. These allow sheaf to find more
// images.
type UserDefinedImage struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	JSONPath   string `json:"jsonPath"`
}

// Validate validates a user defined image.
func (udi UserDefinedImage) Validate() error {
	var apiVersionErr, kindErr, jsonPathErr error

	if udi.APIVersion == "" {
		apiVersionErr = fmt.Errorf("api version is blank")
	}

	if udi.Kind == "" {
		kindErr = fmt.Errorf("kind is blank")
	}

	if udi.JSONPath == "" {
		jsonPathErr = fmt.Errorf("json path is blank")
	} else {
		j := jsonpath.New("parser")
		if err := j.Parse(udi.JSONPath); err != nil {
			jsonPathErr = fmt.Errorf("unable to parse json path %q: %w", udi.JSONPath, err)
		}
	}

	return multierr.Combine(apiVersionErr, kindErr, jsonPathErr)
}

// UserDefinedImageKey is a key describing a UserDefinedImage.
type UserDefinedImageKey struct {
	APIVersion string
	Kind       string
}

type udiMap map[UserDefinedImageKey]UserDefinedImage

func updateUDI(config BundleConfig, fn func(udiMap)) BundleConfig {
	m := udiMap{}
	for _, cur := range config.GetUserDefinedImages() {
		key := UserDefinedImageKey{
			APIVersion: cur.APIVersion,
			Kind:       cur.Kind,
		}

		m[key] = cur
	}

	fn(m)

	list := convertUDIMapToList(m)

	config.SetUserDefinedImages(list)
	return config
}

func convertUDIMapToList(m udiMap) []UserDefinedImage {
	var keys []UserDefinedImageKey
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].APIVersion < keys[j].APIVersion {
			return true
		}
		if keys[i].APIVersion > keys[j].APIVersion {
			return false
		}
		return keys[i].Kind < keys[j].Kind
	})

	var list []UserDefinedImage
	for _, key := range keys {
		list = append(list, m[key])
	}

	return list
}
