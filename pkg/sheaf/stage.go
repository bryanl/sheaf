/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// StageConfig is configuration for Relocate.
type StageConfig struct {
	ArchivePath    string
	RegistryPrefix string
	BundleFactory  BundleFactory
	ImageStager    ImageRelocator
	Archiver       Archiver
}

// Stage stages an archive in a registry.
func Stage(config StageConfig) error {
	fmt.Printf("Relocating images in %s\n\n", config.ArchivePath)

	unpackDir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	defer func() {
		if rErr := os.RemoveAll(unpackDir); rErr != nil {
			log.Printf("remove temporary fs path %q: %v", unpackDir, rErr)
		}
	}()

	if err := config.Archiver.Unarchive(config.ArchivePath, unpackDir); err != nil {
		return fmt.Errorf("unpack fs: %w", err)
	}

	bundle, err := config.BundleFactory(unpackDir)
	if err != nil {
		return fmt.Errorf("open fs: %w", err)
	}

	fmt.Println("Locating images in archive")
	list, err := bundle.Images()
	if err != nil {
		return fmt.Errorf("load images from fs: %w", err)
	}

	fmt.Println("Moving images to new location")
	if err := config.ImageStager.Relocate(unpackDir, config.RegistryPrefix, list.Slice()); err != nil {
		return fmt.Errorf("stage images: %w", err)
	}

	return nil
}
