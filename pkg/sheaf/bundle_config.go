/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"

	"go.uber.org/multierr"
	"k8s.io/client-go/util/jsonpath"

	"github.com/bryanl/sheaf/pkg/images"
)

const (
	// BundleConfigFilename is the filename for a fs config.
	BundleConfigFilename = "bundle.json"

	bundleConfigDefaultVersion = "0.1.0"
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
	APIVersion string               `json:"apiVersion"`
	Kind       string               `json:"kind"`
	JSONPath   string               `json:"jsonPath"`
	Type       UserDefinedImageType `json:"type"`
}

// Validate validates a user defined image.
func (udi UserDefinedImage) Validate() error {
	var apiVersionErr, kindErr, jsonPathErr, typeErr error

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

	if udi.Type == "" {
		typeErr = fmt.Errorf("type is blank")
	} else if !udiContains(udi.Type, []UserDefinedImageType{SingleResult, MultiResult}) {
		typeErr = fmt.Errorf("unknown type %s", udi.Type)
	}

	return multierr.Combine(apiVersionErr, kindErr, jsonPathErr, typeErr)
}

func udiContains(udi UserDefinedImageType, list []UserDefinedImageType) bool {
	for i := range list {
		if udi == list[i] {
			return true
		}
	}

	return false
}

// UserDefinedImageKey is a key describing a UserDefinedImage.
type UserDefinedImageKey struct {
	APIVersion string
	Kind       string
}

// BundleConfig is a fs configuration.
type BundleConfig struct {
	// SchemaVersion is the version of the schema this fs uses.
	SchemaVersion string `json:"schemaVersion"`
	// Name is the name of the fs.
	Name string `json:"name"`
	// Version is the version of the fs.
	Version string `json:"version"`
	// Images is a set of images required by the fs.
	Images *images.Set `json:"images,omitempty"`
	// UserDefinedImages is a list of user defined image locations.
	UserDefinedImages []UserDefinedImage `json:"userDefinedImages,omitempty"`
}

// NewBundleConfig creates a BundleConfig.
func NewBundleConfig(name, version string) BundleConfig {
	if version == "" {
		version = bundleConfigDefaultVersion
	}

	return BundleConfig{
		Name:          name,
		Version:       version,
		SchemaVersion: "v1alpha1",
	}
}

// BundleConfigWriter writes a bundle config.
type BundleConfigWriter func(Bundle, BundleConfig) error
