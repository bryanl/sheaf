/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"

	"github.com/bryanl/sheaf/internal/goutil"
)

// ConfigPull pulls a bundle from a registry.
func ConfigPull(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	if opts.destination == "" {
		return fmt.Errorf("destination is required")
	}

	if opts.reference == "" {
		return fmt.Errorf("reference is required")
	}

	opts.reporter.Headerf("Pull bundle from registry")

	opts.reporter.Reportf("Pulling image %s from registry", opts.reference)

	image, err := opts.imageReader.Read(opts.reference)
	if err != nil {
		return fmt.Errorf("read image with reference %s: %w", opts.reference, err)
	}

	opts.reporter.Reportf("Extracting layers from image")
	layers, err := image.Layers()
	if err != nil {
		return fmt.Errorf("unable to read layers from image: %w", err)
	}

	if len(layers) != 1 {
		return fmt.Errorf("invalid image format: expected 1 layer; got %d layers", len(layers))
	}

	rc, err := layers[0].Compressed()
	if err != nil {
		return fmt.Errorf("get layer: %w", err)
	}

	defer goutil.Close(rc)

	opts.reporter.Reportf("Unarchiving image to %s", opts.destination)

	if err := opts.archiver.Unarchive(rc, opts.destination); err != nil {
		return fmt.Errorf("unarchive image: %w", err)
	}

	return nil
}
