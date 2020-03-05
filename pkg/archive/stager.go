/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/fs"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// StagerOption is a functional option for Stager.
type StagerOption func(s Stager) Stager

// StagerOptionReporter sets the reporter.
func StagerOptionReporter(r reporter.Reporter) StagerOption {
	return func(s Stager) Stager {
		s.reporter = r
		return s
	}
}

// StagerOptionImageRelocator sets the image relocator.
func StagerOptionImageRelocator(ir sheaf.ImageRelocator) StagerOption {
	return func(s Stager) Stager {
		s.imageRelocator = ir
		return s
	}
}

// StagerOptionArchiver sets the archiver.
func StagerOptionArchiver(a sheaf.Archiver) StagerOption {
	return func(s Stager) Stager {
		s.archiver = a
		return s
	}
}

// StagerOptionBundleFactory sets the bundle factory.
func StagerOptionBundleFactory(bf sheaf.BundleFactory) StagerOption {
	return func(s Stager) Stager {
		s.bundleFactory = bf
		return s
	}
}

// StagerOptionBundleReporter sets the bundle reporter.
func StagerOptionBundleReporter(r reporter.Reporter) StagerOption {
	return func(s Stager) Stager {
		s.reporter = r
		return s
	}
}

// Stager stages an archive.
type Stager struct {
	archiver       sheaf.Archiver
	bundleFactory  sheaf.BundleFactory
	imageRelocator sheaf.ImageRelocator
	reporter       reporter.Reporter
}

var _ sheaf.Stager = &Stager{}

// NewStager creates an instance of Stager.
func NewStager(options ...StagerOption) *Stager {
	stager := Stager{
		archiver:       archiver.Default,
		bundleFactory:  fs.DefaultBundleFactory,
		imageRelocator: fs.NewImageRelocator(),
		reporter:       reporter.Default,
	}

	for _, option := range options {
		stager = option(stager)
	}

	return &stager
}

// Stage stages an archive to a registry given a prefix.
func (s Stager) Stage(archiveURI, registryPrefix string) error {
	s.reporter.Headerf("Relocating images in %s", archiveURI)

	unpackDir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	defer func() {
		if rErr := os.RemoveAll(unpackDir); rErr != nil {
			log.Printf("remove temporary fs path %q: %v", unpackDir, rErr)
		}
	}()

	if err := s.archiver.Unarchive(archiveURI, unpackDir); err != nil {
		return fmt.Errorf("unpack fs: %w", err)
	}

	bundle, err := s.bundleFactory(unpackDir)
	if err != nil {
		return fmt.Errorf("open bundle: %w", err)
	}

	s.reporter.Report("Locating images in archive")
	list, err := bundle.Images()
	if err != nil {
		return fmt.Errorf("load images from fs: %w", err)
	}

	s.reporter.Header("Moving images to new location")
	if err := s.imageRelocator.Relocate(unpackDir, registryPrefix, list.Slice()); err != nil {
		return fmt.Errorf("stage images: %w", err)
	}

	return nil
}
