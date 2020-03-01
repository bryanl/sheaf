/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"log"
	"os"
)

// PackConfig is configuration for Pack.
type PackConfig struct {
	Packer        Packer
	BundleURI     string
	BundleFactory BundleFactory
}

// Pack packs a bundle.
func Pack(config PackConfig) error {
	bundle, err := config.BundleFactory(config.BundleURI)
	if err != nil {
		return fmt.Errorf("open bundle: %w", err)
	}

	bundleConfig := bundle.Config()
	dest := bundleConfig.Filename(".")

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create archive file: %w", err)
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("close archive file: %v", err)
		}
	}()

	if err := config.Packer.Pack(bundle, f); err != nil {
		return fmt.Errorf("pack archive: %w", err)
	}

	return nil
}
