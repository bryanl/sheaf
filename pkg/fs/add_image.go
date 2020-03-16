/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"github.com/bryanl/sheaf/pkg/sheaf"
	"github.com/pivotal/image-relocation/pkg/images"
)

// ImageAdderOption is a functional option for ImageAdder.
type ImageAdderOption func(ia ImageAdder) ImageAdder

// ImageAdder adds images to a fs.
type ImageAdder struct {
	bundlePath         string
	bundleFactory      sheaf.BundleFactory
	bundleConfigWriter sheaf.BundleConfigWriter
}

// NewImageAdder creates an instance of ImageAdder.
func NewImageAdder(bundlePath string, options ...ImageAdderOption) (*ImageAdder, error) {
	ia := ImageAdder{
		bundlePath:         bundlePath,
		bundleFactory:      DefaultBundleFactory,
		bundleConfigWriter: DefaultBundleConfigWriter,
	}

	for _, option := range options {
		ia = option(ia)
	}

	return &ia, nil
}

// Add adds a list of images to the fs manifest.
func (ia *ImageAdder) Add(imageList ...string) error {
	b, err := ia.bundleFactory(ia.bundlePath)
	if err != nil {
		return err
	}

	bc := b.Config()

	im, err := images.New(imageList...)
	if err != nil {
		return err
	}

	cur := images.Empty
	if bc.Images != nil {
		cur = *bc.Images
	}

	n := cur.Union(im)
	bc.Images = &n

	return ia.bundleConfigWriter(b, bc)
}
