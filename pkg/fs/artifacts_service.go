/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"io/ioutil"
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// ArtifactsService is a service for interacting with artifacts in a fs.
type ArtifactsService struct {
	bundle sheaf.Bundle
}

var _ sheaf.ArtifactsService = &ArtifactsService{}

// NewArtifactsService creates an instance of ArtifactsService.
func NewArtifactsService(bundle sheaf.Bundle) *ArtifactsService {
	s := ArtifactsService{
		bundle: bundle,
	}

	return &s
}

// Index returns the contents of an artifact layout index.json as bytes.
func (s *ArtifactsService) Index() ([]byte, error) {
	layoutPath := filepath.Join(s.bundle.Path(), "artifacts", "layout")
	index := filepath.Join(layoutPath, "index.json")

	return ioutil.ReadFile(index)
}

// Image is an image service for interacting with images in the artifacts.
func (s *ArtifactsService) Image() sheaf.ImageService {
	return NewImageService(s)
}
