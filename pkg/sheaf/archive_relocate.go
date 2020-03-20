/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
)

// ArchiveRelocate relocates images to an archive to a registry.
func ArchiveRelocate(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	opts.reporter.Headerf("Relocating images in %s", opts.archive)

	return withExplodedArchive(opts, func(b Bundle) error {
		opts.reporter.Report("Locating images in archive")
		list, err := b.Images()
		if err != nil {
			return fmt.Errorf("load images from fs: %w", err)
		}

		opts.reporter.Header("Moving images to new location")

		if err := opts.imageRelocator.Relocate(b.Path(), opts.repositoryPrefix, list.Slice(), opts.imageWriter); err != nil {
			return fmt.Errorf("stage images: %w", err)
		}

		return nil
	})
}
