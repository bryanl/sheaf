/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/codec"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// PackerOption is a functional option for configuring Packer.
type PackerOption func(p Packer) Packer

// PackerLayoutFactory sets the layout factory for Packer.
func PackerLayoutFactory(lf LayoutFactory) PackerOption {
	return func(p Packer) Packer {
		p.layoutFactory = lf
		return p
	}
}

// PackerArchiver sets the archiver for Packer.
func PackerArchiver(a sheaf.Archiver) PackerOption {
	return func(p Packer) Packer {
		p.archiver = a
		return p
	}
}

// PackerReporter sets the reporter for Packer.
func PackerReporter(r reporter.Reporter) PackerOption {
	return func(p Packer) Packer {
		p.reporter = r
		return p
	}
}

// Packer packs a fs into an archive.
type Packer struct {
	codec         sheaf.Encoder
	archiver      sheaf.Archiver
	layoutFactory LayoutFactory
	reporter      reporter.Reporter
}

var _ sheaf.Packer = &Packer{}

// NewPacker creates an instance of Packer.
func NewPacker(options ...PackerOption) *Packer {
	p := Packer{
		codec:         codec.DefaultEncoder,
		archiver:      archiver.Default,
		layoutFactory: DefaultLayoutFactory(),
		reporter:      reporter.Default,
	}

	for _, option := range options {
		p = option(p)
	}

	return &p
}

// Pack runs the pack operation.
func (p Packer) Pack(b sheaf.Bundle, w io.Writer) error {
	dir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}

	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			log.Printf("unable to remove temporary directory: %v", err)
		}
	}()

	if err := p.stageBundleConfig(dir, b); err != nil {
		return fmt.Errorf("stage fs config: %w", err)
	}

	if err := p.stageManifests(dir, b); err != nil {
		return fmt.Errorf("stage manifests: %w", err)
	}

	if err := p.stageImages(dir, b); err != nil {
		return fmt.Errorf("stage images")
	}

	p.reporter.Header("Creating archive")
	if err := p.archiver.Archive(dir, w); err != nil {
		return fmt.Errorf("create packed archive: %w", err)
	}

	return nil
}

func (p Packer) stageImages(dir string, b sheaf.Bundle) error {
	p.reporter.Header("Staging images")

	layout, err := p.layoutFactory(dir)
	if err != nil {
		return fmt.Errorf("create layout manager: %w", err)
	}

	imageList, err := b.Images()
	if err != nil {
		return fmt.Errorf("get images from fs: %w", err)
	}

	for _, imageName := range imageList.Slice() {
		p.reporter.Reportf("adding %s to layout\n", imageName.String())
		if _, err := layout.Add(imageName); err != nil {
			return fmt.Errorf("add ref %s to image layout: %w", imageName, err)
		}
	}
	return nil
}

func (p Packer) stageManifests(dir string, b sheaf.Bundle) error {
	p.reporter.Header("Staging manifests")

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

func (p Packer) stageBundleConfig(dir string, b sheaf.Bundle) error {
	p.reporter.Header("Staging bundle configuration")

	bundleConfigPath := filepath.Join(dir, sheaf.BundleConfigFilename)
	bundleConfigData, err := p.codec.Encode(b.Config())
	if err != nil {
		return fmt.Errorf("encode fs config: %w", err)
	}
	if err := ioutil.WriteFile(bundleConfigPath, bundleConfigData, 0600); err != nil {
		return fmt.Errorf("write fs config: %w", err)
	}
	return nil
}
