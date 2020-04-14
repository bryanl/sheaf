/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"io"
)

//go:generate mockgen -destination=../mocks/mock_bundle_config.go -package mocks github.com/bryanl/sheaf/pkg/sheaf BundleConfig
//go:generate mockgen -destination=../mocks/mock_bundle_config_codec.go -package mocks github.com/bryanl/sheaf/pkg/sheaf BundleConfigCodec
//go:generate mockgen -destination=../mocks/mock_bundle_config_writer.go -package mocks github.com/bryanl/sheaf/pkg/sheaf BundleConfigWriter

const (
	// BundleConfigFilename is the filename for a fs config.
	BundleConfigFilename = "bundle.json"

	// BundleConfigDefaultVersion is the default version a bundle.
	BundleConfigDefaultVersion = "0.1.0"
)

// BundleConfig is a bundle configuration interface.
type BundleConfig interface {
	GetSchemaVersion() string
	SetSchemaVersion(string)
	GetName() string
	SetName(string)
	GetVersion() string
	SetVersion(string)
	GetUserDefinedImages() []UserDefinedImage
	SetUserDefinedImages([]UserDefinedImage)
}

// BundleConfigWriter writes a bundle config.
type BundleConfigWriter interface {
	Write(Bundle, BundleConfig) error
}

// BundleConfigCodec is a codec for bundle config.
type BundleConfigCodec interface {
	// Encode encodes a bundle config to a writer.
	Encode(w io.Writer, bc BundleConfig) error
	// Decode decodes a bundle config from a reader.
	Decode(r io.Reader) (BundleConfig, error)
}
