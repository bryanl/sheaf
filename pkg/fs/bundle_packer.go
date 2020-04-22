/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/internal/goutil"
	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// BundlePackerOption is a functional option for configuring BundlePacker.
type BundlePackerOption func(bp *BundlePacker)

// BundlePacker packs bundles that live on a filesystem.
type BundlePacker struct {
	reporter           reporter.Reporter
	archiver           sheaf.Archiver
	layoutFactory      LayoutFactory
	bundleConfigWriter sheaf.BundleConfigWriter
}

var _ sheaf.BundlePacker = &BundlePacker{}

// NewBundlePacker creates an instance of BundlePacker.
func NewBundlePacker(options ...BundlePackerOption) *BundlePacker {
	bp := BundlePacker{
		reporter:           reporter.New(),
		archiver:           archiver.New(),
		layoutFactory:      DefaultLayoutFactory(),
		bundleConfigWriter: NewBundleConfigWriter(),
	}

	for _, option := range options {
		option(&bp)
	}

	return &bp
}

// Pack packs a bundle to a filesystem destination.
func (bp BundlePacker) Pack(b sheaf.Bundle, dest string, force bool) error {
	bundleConfig := b.Config()

	filename := fmt.Sprintf("%s-%s.tgz", bundleConfig.GetName(), bundleConfig.GetVersion())

	dest = filepath.Join(dest, filename)
	if _, err := os.Stat(dest); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

	} else {
		if force {
			bp.reporter.Reportf("Removing existing file at %s", dest)
			if err = os.RemoveAll(dest); err != nil {
				return fmt.Errorf("unable to remove destination: %w", err)
			}
		} else {
			return fmt.Errorf("destination %s exists", dest)
		}
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create archive file: %w", err)
	}

	defer goutil.Close(f)

	dir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}

	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			log.Printf("unable to remove temporary directory: %v", err)
		}
	}()

	if err := bp.stageBundleConfig(dir, b); err != nil {
		return fmt.Errorf("stage bundle config: %w", err)
	}

	if err := bp.stageManifests(dir, b); err != nil {
		return fmt.Errorf("stage manifests: %w", err)
	}

	if err := bp.stageImages(dir, b); err != nil {
		return fmt.Errorf("stage images: %w", err)
	}

	bp.reporter.Headerf("Creating archive: %s", dest)
	if err := bp.archiver.Archive(dir, f); err != nil {
		return fmt.Errorf("create packed archive: %w", err)
	}

	return nil
}

func (bp BundlePacker) stageImages(dir string, b sheaf.Bundle) error {
	bp.reporter.Header("Staging images")

	layout, err := bp.layoutFactory(dir)
	if err != nil {
		return fmt.Errorf("create layout manager: %w", err)
	}

	imageList, err := b.Images()
	if err != nil {
		return fmt.Errorf("get images from fs: %w", err)
	}

	for _, imageName := range imageList.Slice() {
		bp.reporter.Reportf("adding %s to layout\n", imageName.String())
		if _, err := layout.Add(imageName); err != nil {
			return fmt.Errorf("add ref %s to image layout: %w", imageName, err)
		}
	}
	return nil
}

func (bp BundlePacker) stageManifests(dir string, b sheaf.Bundle) error {
	bp.reporter.Header("Staging manifests")

	manifestsDest := filepath.Join(dir, "app", "manifests")
	if err := os.MkdirAll(manifestsDest, 0700); err != nil {
		return fmt.Errorf("create manifests directory: %w", err)
	}

	m, err := b.Manifests()
	if err != nil {
		return err
	}

	bundleManifests, err := m.List()
	if err != nil {
		return fmt.Errorf("get manifest paths: %w", err)
	}

	for _, bundleManifest := range bundleManifests {
		name := filepath.Base(bundleManifest.ID)

		dest := filepath.Join(manifestsDest, name)
		if err := ioutil.WriteFile(dest, bundleManifest.Data, 0600); err != nil {
			return fmt.Errorf("stage manifest %q: %w", bundleManifest.ID, err)
		}
	}
	return nil
}

func (bp BundlePacker) stageBundleConfig(dir string, b sheaf.Bundle) error {
	bp.reporter.Header("Staging bundle configuration")

	if _, err := b.Copy(dir); err != nil {
		return fmt.Errorf("duplicate bundle to tempoary directory: %w", err)
	}

	return nil
}
