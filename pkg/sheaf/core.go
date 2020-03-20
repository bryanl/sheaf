/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"github.com/pivotal/image-relocation/pkg/image"
)

//go:generate mockgen -destination=../mocks/mock_artifacts_service.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ArtifactsService
//go:generate mockgen -destination=../mocks/mock_image_service.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ImageService
//go:generate mockgen -destination=../mocks/mock_image_relocator.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ImageRelocator

// BundleImage is an image in a fs.
type BundleImage struct {
	Name      string `json:"name"`
	Digest    string `json:"digest"`
	MediaType string `json:"mediaType"`
}

// ImageService returns a list of fs artifact images.
type ImageService interface {
	List() ([]BundleImage, error)
}

// ArtifactsService interacts with fs artifacts.
type ArtifactsService interface {
	Index() ([]byte, error)
	Image() ImageService
}

// ImageRelocator relocates an images to another registry.
type ImageRelocator interface {
	Relocate(rootPath, prefix string, images []image.Name, iw ImageWriter) error
}
