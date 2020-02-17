/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
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

	"github.com/pivotal/image-relocation/pkg/image"
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
	images, err := imagesFromStrings(imageStrs)
	if err != nil {
		return err
	}

	bc, bcPath, err := loadBundleConfig(ia.BundlePath)
	if err != nil {
		return err
	}

	bundleImages, err := imagesFromStrings(bc.Images)
	if err != nil {
		return err
	}
	bc.Images = imageStrings(union(bundleImages, images))

	return StoreBundleConfig(bc, bcPath)
}

func union(a []image.Name, b []image.Name) []image.Name {
	uniq := map[image.Name]struct{}{}

	for _, i := range a {
		uniq[i] = struct{}{}
	}

	for _, i := range b {
		uniq[i] = struct{}{}
	}

	imgs := []image.Name{}
	for i := range uniq {
		imgs = append(imgs, i)
	}

	// Enforce a deterministic ordering, e.g for testing and repeatable building.
	sortImages(imgs)

	return imgs
}
