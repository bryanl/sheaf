/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"log"
)

// PackConfig is configuration for Pack.
type PackConfig struct {
	Path string
}

// Pack packs a bundle.
func Pack(config PackConfig) error {
	bundle, err := OpenBundle(config.Path)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	defer func() {
		if cErr := bundle.Close(); cErr != nil {
			log.Printf("unable to close bundle: %v", err)
		}
	}()

	imgs, err := bundle.Images()
	if err != nil {
		return fmt.Errorf("collect images from manifest: %w", err)
	}

	for _, ref := range imgs.Slice() {
		fmt.Printf("Adding %s to bundle\n", ref)
		if _, err := bundle.Layout.Add(ref); err != nil {
			return fmt.Errorf("add %s: %w", ref, err)
		}
	}

	if err := bundle.Write(); err != nil {
		return fmt.Errorf("write bundle archive: %w", err)
	}

	return nil
}
