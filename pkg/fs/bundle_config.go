/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/pivotal/image-relocation/pkg/images"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// BundleConfigCodec is a codec for encoding/decoding bundle configs.
type BundleConfigCodec struct {
}

var _ sheaf.BundleConfigCodec = &BundleConfigCodec{}

// NewBundleConfigCodec creates an instance of BundleConfigCodec.
func NewBundleConfigCodec() *BundleConfigCodec {
	bcc := BundleConfigCodec{}

	return &bcc
}

// Encode encodes a bundle config.
func (b BundleConfigCodec) Encode(w io.Writer, bc sheaf.BundleConfig) error {
	var imageList []string
	if list := bc.GetImages(); list != nil {
		imageList = list.Strings()
	}

	bcf := bundleConfigFormat{
		SchemaVersion:     bc.GetSchemaVersion(),
		Name:              bc.GetName(),
		Version:           bc.GetVersion(),
		Images:            imageList,
		UserDefinedImages: bc.GetUserDefinedImages(),
	}

	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(&bcf)
}

// Decode decodes b bundle config.
func (b BundleConfigCodec) Decode(r io.Reader) (sheaf.BundleConfig, error) {
	var bcf bundleConfigFormat
	if err := json.NewDecoder(r).Decode(&bcf); err != nil {
		return nil, err
	}

	imageList, err := images.New(bcf.Images...)
	if err != nil {
		return nil, fmt.Errorf("images are invalid: %w", err)
	}

	bc := BundleConfig{
		schemaVersion:     bcf.SchemaVersion,
		name:              bcf.Name,
		version:           bcf.Version,
		images:            &imageList,
		userDefinedImages: bcf.UserDefinedImages,
	}

	return &bc, nil
}

type bundleConfigFormat struct {
	// SchemaVersion is the version of the schema this fs uses.
	SchemaVersion string `json:"schemaVersion"`
	// Name is the name of the fs.
	Name string `json:"name"`
	// Version is the version of the fs.
	Version string `json:"version"`
	// Images is a set of images required by the fs.
	Images []string `json:"images,omitempty"`
	// UserDefinedImages is a list of user defined image locations.
	UserDefinedImages []sheaf.UserDefinedImage `json:"userDefinedImages,omitempty"`
}

// BundleConfig is a bundle configuration.
type BundleConfig struct {
	// SchemaVersion is the version of the schema this fs uses.
	schemaVersion string
	// Name is the name of the fs.
	name string
	// Version is the version of the fs.
	version string
	// Images is a set of images required by the fs.
	images *images.Set
	// UserDefinedImages is a list of user defined image locations.
	userDefinedImages []sheaf.UserDefinedImage
}

var _ sheaf.BundleConfig = &BundleConfig{}

// GetSchemaVersion returns the bundle config's schema version.
func (b BundleConfig) GetSchemaVersion() string {
	return b.schemaVersion
}

// SetSchemaVersion sets the bundle config's schema version.
func (b *BundleConfig) SetSchemaVersion(schemaVersion string) {
	b.schemaVersion = schemaVersion
}

// GetName returns the bundle config's name.
func (b BundleConfig) GetName() string {
	return b.name
}

// SetName sets the bundle config's name.
func (b *BundleConfig) SetName(name string) {
	b.name = name
}

// GetVersion returns the bundle config's version.
func (b BundleConfig) GetVersion() string {
	return b.version
}

// SetVersion sets the bundle config's version.
func (b *BundleConfig) SetVersion(version string) {
	b.version = version
}

// GetImages returns the bundle config's images.
func (b BundleConfig) GetImages() *images.Set {
	return b.images
}

// SetImages sets the bundle config's images.
func (b *BundleConfig) SetImages(imageSet *images.Set) {
	b.images = imageSet
}

// GetUserDefinedImages returns the bundle config's user defined images.
func (b BundleConfig) GetUserDefinedImages() []sheaf.UserDefinedImage {
	return b.userDefinedImages
}

// SetUserDefinedImages sets the bundle config's user defined images.
func (b *BundleConfig) SetUserDefinedImages(userDefinedImages []sheaf.UserDefinedImage) {
	b.userDefinedImages = userDefinedImages
}

// NewBundleConfig creates a BundleConfig.
func NewBundleConfig(name, version string) *BundleConfig {
	if version == "" {
		version = sheaf.BundleConfigDefaultVersion
	}

	return &BundleConfig{
		name:          name,
		version:       version,
		schemaVersion: "v1alpha1",
	}
}
