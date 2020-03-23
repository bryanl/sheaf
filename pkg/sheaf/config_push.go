/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
)

// ConfigPush pushes a bundle to a registry.
func ConfigPush(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	if opts.reference == "" {
		return fmt.Errorf("reference is required")
	}

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	opts.reporter.Headerf("Push bundle to %s", opts.reference)

	image, err := opts.bundleImager.CreateImage(b)
	if err != nil {
		return fmt.Errorf("create image from bundle: %w", err)
	}

	if err := opts.imageWriter.Write(opts.reference, image); err != nil {
		return fmt.Errorf("write image to registry: %w", err)
	}

	return nil
}
