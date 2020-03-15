/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// PackConfig is configuration for Pack.
type PackConfig struct {
	Packer        sheaf.Packer
	BundleURI     string
	BundleFactory sheaf.BundleFactory
	Dest          string
	Force         bool
}

// Pack packs a fs.
func Pack(config PackConfig) error {
	bundle, err := config.BundleFactory(config.BundleURI)
	if err != nil {
		return fmt.Errorf("open bundle at %s: %w", config.BundleURI, err)
	}

	bundleConfig := bundle.Config()

	filename := fmt.Sprintf("%s-%s.tgz", bundleConfig.Name, bundleConfig.Version)

	dest := filepath.Join(config.Dest, filename)
	if _, err := os.Stat(dest); err == nil && !config.Force {
		return fmt.Errorf("unable to create archive %q because a file already exists", dest)
	} else if !os.IsNotExist(err) {
		return err
	}

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
