/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"strings"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/pathmapping"

	"github.com/bryanl/sheaf/pkg/manifest"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// ImageReplacer replaces images in manifests on a filesystem.
type ImageReplacer struct{}

var _ sheaf.ImageReplacer = &ImageReplacer{}

// NewImageReplacer creates an instance of ImageReplacer.
func NewImageReplacer() *ImageReplacer {
	ir := ImageReplacer{}

	return &ir
}

// Replace replaces container images found in a bundle manifest.
func (i ImageReplacer) Replace(m sheaf.BundleManifest, config sheaf.BundleConfig, prefix string) ([]byte, error) {
	data := m.Data

	if prefix != "" {
		imageSet, err := manifest.ContainerImagesFromBytes(data, config.GetUserDefinedImages())
		if err != nil {
			return nil, fmt.Errorf("container images from manifest: %w", err)
		}

		imageMap := make(map[image.Name]image.Name)
		for _, img := range imageSet.Slice() {
			newImageName, err := pathmapping.FlattenRepoPathPreserveTagDigest(prefix, img)
			if err != nil {
				return nil, fmt.Errorf("create flatten image name from %s with prefix %q: %w", img.String(), prefix, err)
			}
			imageMap[img] = newImageName
		}

		data = replaceImage(data, imageMap)
	}

	return data, nil
}

func replaceImage(manifest []byte, imageMap map[image.Name]image.Name) []byte {
	var replacements []string
	for oldImage, newImage := range imageMap {
		for _, oi := range oldImage.Synonyms() {
			replacements = append(replacements, oi.String(), newImage.String())
		}
	}
	return []byte(strings.NewReplacer(replacements...).Replace(string(manifest)))
}
