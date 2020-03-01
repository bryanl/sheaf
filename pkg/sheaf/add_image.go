/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"

	"github.com/bryanl/sheaf/pkg/images"
)

// ImageAdder adds images to a fs.
type ImageAdder struct {
	// BundlePath is the path to the fs.
	BundlePath string
}

// NewImageAdder creates an instance of ImageAdder.
func NewImageAdder(bundlePath string) (*ImageAdder, error) {
	if err := ensureBundlePath(bundlePath); err != nil {
		return nil, fmt.Errorf("ensure fs path %q exists: %w", bundlePath, err)
	}

	return &ImageAdder{
		BundlePath: bundlePath,
	}, nil
}

// Add adds a list of images to the fs manifest.
func (ia *ImageAdder) Add(imageStrs []string) error {
	im, err := images.New(imageStrs)
	if err != nil {
		return err
	}

	bc, bcPath, err := loadBundleConfig(ia.BundlePath)
	if err != nil {
		return err
	}

	bc.Images = bc.Images.Union(im)

	return StoreBundleConfig(bc, bcPath)
}
