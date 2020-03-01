/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"sort"

	ociv1 "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/bryanl/sheaf/pkg/codec"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// ImageServiceOption is a functional option for configuration ImageService.
type ImageServiceOption func(is ImageService) ImageService

// ImageServiceDecoder sets the decoder for ImageService.
func ImageServiceDecoder(d sheaf.Decoder) ImageServiceOption {
	return func(is ImageService) ImageService {
		is.Decoder = d
		return is
	}
}

// ImageService uses an ArtifactsService to interact with images.
type ImageService struct {
	ArtifactsService sheaf.ArtifactsService
	Decoder          sheaf.Decoder
}

var _ sheaf.ImageService = &ImageService{}

// NewImageService creates an instance of ImageService.
func NewImageService(artifactsService sheaf.ArtifactsService, options ...ImageServiceOption) *ImageService {
	is := ImageService{
		ArtifactsService: artifactsService,
	}

	for _, option := range options {
		is = option(is)
	}

	if is.Decoder == nil {
		is.Decoder = codec.DefaultDecoder
	}
	return &is
}

// List lists images from an OCI index.
func (s *ImageService) List() ([]sheaf.BundleImage, error) {
	data, err := s.ArtifactsService.Index()
	if err != nil {
		return nil, err
	}

	var index ociv1.Index
	if err := s.Decoder.Decode(data, &index); err != nil {
		return nil, err
	}

	var list []sheaf.BundleImage

	for _, manifest := range index.Manifests {
		refName, ok := manifest.Annotations[ociv1.AnnotationRefName]
		if !ok {
			continue
		}

		image := sheaf.BundleImage{
			Name:      refName,
			Digest:    manifest.Digest.String(),
			MediaType: manifest.MediaType,
		}

		list = append(list, image)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Name < list[j].Name
	})

	return list, nil
}
