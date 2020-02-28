/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

//go:generate mockgen -destination=../mocks/mock_bundle_service.go -package mocks github.com/bryanl/sheaf/pkg/sheaf BundleService
//go:generate mockgen -destination=../mocks/mock_artifacts_service.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ArtifactsService
//go:generate mockgen -destination=../mocks/mock_decoder.go -package mocks github.com/bryanl/sheaf/pkg/sheaf Decoder
//go:generate mockgen -destination=../mocks/mock_image_service.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ImageService

// BundleImage is an image in a bundle.
type BundleImage struct {
	Name      string `json:"name"`
	Digest    string `json:"digest"`
	MediaType string `json:"mediaType"`
}

// ImageService returns a list of bundle artifact images.
type ImageService interface {
	List() ([]BundleImage, error)
}

// ArtifactsService interacts with bundle artifacts.
type ArtifactsService interface {
	Index() ([]byte, error)
	Image() ImageService
}

// Decoder decodes bytes into a value.
type Decoder interface {
	Decode([]byte, interface{}) error
}

// Encoder encodes a value into bytes.
type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

// Codec combines Decoder and Encoder
type Codec interface {
	Decoder
	Encoder
}

// BundleService manages bundles.
type BundleService interface {
	Codec() Codec
	Path() string
	Config() BundleConfig
	Artifacts() ArtifactsService
}

// Archiver manages archives.
type Archiver interface {
	Unarchive(src, dest string) error
}
