/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"path/filepath"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/pathmapping"

	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// ImageRelocatorOption is a functional option for configuring ImageRelocator.
type ImageRelocatorOption func(is ImageRelocator) ImageRelocator

// ImageRelocatorLayoutFactory configuration the layout factory.
func ImageRelocatorLayoutFactory(lf LayoutFactory) ImageRelocatorOption {
	return func(is ImageRelocator) ImageRelocator {
		is.layoutFactory = lf
		return is
	}
}

// ImageRelocatorDryRun configures ImageRelocator to do a dry run.
func ImageRelocatorDryRun(dryRun bool) ImageRelocatorOption {
	return func(is ImageRelocator) ImageRelocator {
		is.dryRun = dryRun
		return is
	}
}

// ImageRelocator relocates images to a registry.
type ImageRelocator struct {
	layoutFactory LayoutFactory
	reporter      reporter.Reporter
	dryRun        bool
}

var _ sheaf.ImageRelocator = &ImageRelocator{}

// NewImageRelocator creates an instance of ImageRelocator.
func NewImageRelocator(options ...ImageRelocatorOption) *ImageRelocator {
	is := ImageRelocator{
		layoutFactory: DefaultLayoutFactory(),
		reporter:      reporter.Default,
	}

	for _, option := range options {
		is = option(is)
	}

	return &is
}

// Relocate relocates images to a registry given a prefix.
func (i ImageRelocator) Relocate(rootPath, prefix string, images []image.Name, iw sheaf.ImageWriter) error {
	layoutPath := filepath.Join(rootPath)
	l, err := i.layoutFactory(layoutPath)
	if err != nil {
		return fmt.Errorf("create layout: %w", err)
	}

	for _, imageName := range images {
		imageDigest, err := l.Find(imageName)
		if err != nil {
			return fmt.Errorf("find image digest for ref %s: %w", imageName.String(), err)
		}

		newImageName, err := pathmapping.FlattenRepoPathPreserveTagDigest(prefix, imageName)
		if err != nil {
			return fmt.Errorf("create relocated image name: %w", err)
		}

		i.printImageRelocation(imageName.String(), newImageName.String(), i.dryRun)
		if i.dryRun {
			continue
		}

		// TODO: need to support insecure images
		if err := l.Push(imageDigest, newImageName); err != nil {
			return fmt.Errorf("push %s: %w", newImageName.String(), err)
		}
	}

	return nil
}

func (i ImageRelocator) printImageRelocation(oldName, newName string, isDryRun bool) {
	var marker string
	if isDryRun {
		marker = " (DRY RUN)"
	}
	i.reporter.Reportf("Relocating image%s\n%s old: %s\n%s new: %s",
		marker, treeItem, oldName, treeItemLast, newName)
}
