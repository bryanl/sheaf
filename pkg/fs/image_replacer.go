/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
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
	if prefix == "" {
		return m.Data, nil
	}

	return manifest.MapContainer(m.Data, config.GetUserDefinedImages(), func(originalImage image.Name) (image.Name, error) {
		return pathmapping.FlattenRepoPathPreserveTagDigest(prefix, originalImage)
	})
}
