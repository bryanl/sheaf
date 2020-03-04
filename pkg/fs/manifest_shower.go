/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/pathmapping"

	"github.com/bryanl/sheaf/pkg/manifest"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// ManifestShowerOption is an option for configuring ManifestShower.
type ManifestShowerOption func(mg ManifestShower) ManifestShower

// ManifestShowerPrefix sets the prefix for relocating images.
func ManifestShowerPrefix(p string) ManifestShowerOption {
	return func(mg ManifestShower) ManifestShower {
		mg.Prefix = p
		return mg
	}
}

// ManifestShowerArchiver sets the archiver for ManifestShower.
func ManifestShowerArchiver(a sheaf.Archiver) ManifestShowerOption {
	return func(mg ManifestShower) ManifestShower {
		mg.Archiver = a
		return mg
	}
}

// ManifestShower generates manifests from a fs archive.
type ManifestShower struct {
	Path     string
	Prefix   string
	Archiver sheaf.Archiver
}

// NewManifestShower creates an instance of ManifestShower.
func NewManifestShower(p string, options ...ManifestShowerOption) *ManifestShower {
	mg := ManifestShower{
		Path: p,
	}

	for _, option := range options {
		mg = option(mg)
	}

	return &mg
}

// Show generates the manifests contained in a fs archive to the supplied writer.
func (mg *ManifestShower) Show(w io.Writer) error {
	bundleDir := mg.Path

	fi, err := os.Stat(bundleDir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		// assume path is an archive
		dir, err := mg.extractArchive(mg.Path)
		if err != nil {
			return err
		}

		bundleDir = dir

		defer func() {
			if rErr := os.RemoveAll(dir); rErr != nil {
				log.Printf("unable to remove temporary directory: %v", err)
			}
		}()
	} else {
		bundleDir, err = locateRootDir(bundleDir)
		if err != nil {
			return err
		}
	}

	manifestsPath := filepath.Join(bundleDir, "app", "manifests")

	b, err := NewBundle(bundleDir)
	if err != nil {
		return err
	}

	config := b.Config()

	entries, err := ioutil.ReadDir(manifestsPath)
	if err != nil {
		return err
	}

	for i := range entries {
		fi := entries[i]
		if fi.IsDir() {
			continue
		}

		manifestPath := filepath.Join(manifestsPath, fi.Name())
		data, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return err
		}

		if mg.Prefix != "" {
			images, err := manifest.ContainerImages(manifestPath, config.UserDefinedImages)
			if err != nil {
				return err
			}

			imageMap := make(map[image.Name]image.Name)
			for _, img := range images.Slice() {
				newImageName, err := pathmapping.FlattenRepoPathPreserveTagDigest(mg.Prefix, img)
				if err != nil {
					return err
				}
				imageMap[img] = newImageName
			}

			data = replaceImage(data, imageMap)
		}

		if i > 0 {
			if _, err := fmt.Fprintln(w, "---"); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(w, string(bytes.TrimSpace(data))); err != nil {
			return err
		}
	}

	return nil
}

func (mg *ManifestShower) extractArchive(archivePath string) (string, error) {
	tmpDir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return "", fmt.Errorf("create temporary directory: %w", err)
	}

	if err := mg.Archiver.Unarchive(archivePath, tmpDir); err != nil {
		if rErr := os.RemoveAll(tmpDir); rErr != nil {
			log.Printf("remove temporary directory: %v", rErr)
		}
		return "", fmt.Errorf("unpack fs: %w", err)
	}

	return tmpDir, nil
}
