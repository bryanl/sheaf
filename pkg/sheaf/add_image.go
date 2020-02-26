/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"fmt"

	"github.com/bryanl/sheaf/pkg/images"
)

// ImageAdder adds images to a bundle.
type ImageAdder struct {
	// BundlePath is the path to the bundle.
	BundlePath string
}

// NewImageAdder creates an instance of ImageAdder.
func NewImageAdder(bundlePath string) (*ImageAdder, error) {
	if err := ensureBundlePath(bundlePath); err != nil {
		return nil, fmt.Errorf("ensure bundle path %q exists: %w", bundlePath, err)
	}

	return &ImageAdder{
		BundlePath: bundlePath,
	}, nil
}

// Add adds a list of images to the bundle manifest.
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
