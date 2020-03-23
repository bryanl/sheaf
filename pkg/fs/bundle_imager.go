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

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// BundleImagerOption is a functional option for configuring BundleImager.
type BundleImagerOption func(bi *BundleImager)

// BundleImagerReporter sets the reporter.
func BundleImagerReporter(r reporter.Reporter) BundleImagerOption {
	return func(bi *BundleImager) {
		bi.reporter = r
	}
}

// BundleImager creates an image from a bundle that lives on a filesystem.
type BundleImager struct {
	archiver sheaf.Archiver
	reporter reporter.Reporter
}

var _ sheaf.BundleImager = &BundleImager{}

// NewBundleImager creates an instance of BundleImager.
func NewBundleImager(options ...BundleImagerOption) *BundleImager {
	bi := BundleImager{
		archiver: archiver.New(),
		reporter: reporter.New(),
	}

	for _, option := range options {
		option(&bi)
	}

	return &bi
}

// CreateImage create an image from a bundle.
func (bi BundleImager) CreateImage(b sheaf.Bundle) (v1.Image, error) {
	archiveBytes, err := bi.createArchive(b.Path())
	if err != nil {
		return nil, fmt.Errorf("create archive: %w", err)
	}

	var layers []mutate.Addendum

	layer, err := bi.createLayer(archiveBytes)
	if err != nil {
		return nil, fmt.Errorf("create bundle layer: %w", err)
	}

	layers = append(layers, layer)

	r := bi.reporter

	r.Report("Create image from layers")
	base := empty.Image
	withConfig, err := mutate.Append(base, layers...)
	if err != nil {
		return nil, fmt.Errorf("append data layer to base: %w", err)
	}

	r.Report("Setting image configuration")
	cfg, err := withConfig.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("get config file: %w", err)
	}

	cfg = cfg.DeepCopy()
	cfg.Author = "github.com/bryanl/sheaf"

	r.Report("Adding image configuration to image")
	image, err := mutate.ConfigFile(withConfig, cfg)
	if err != nil {
		return nil, fmt.Errorf("mutate config file: %w", err)
	}

	return image, nil
}

func (bi *BundleImager) createArchive(bundlePath string) ([]byte, error) {
	bi.reporter.Reportf("Creating archive of configuration and manifests in %s", bundlePath)
	var b bytes.Buffer
	if err := bi.archiver.Archive(bundlePath, &b); err != nil {
		return nil, fmt.Errorf("create archive: %w", err)
	}

	return b.Bytes(), nil
}

func (bi *BundleImager) createLayer(b []byte) (mutate.Addendum, error) {
	bi.reporter.Report("Create layer with bundle")
	dataLayer, err := tarball.LayerFromOpener(func() (closer io.ReadCloser, err error) {
		return ioutil.NopCloser(bytes.NewBuffer(b)), nil
	})
	if err != nil {
		return mutate.Addendum{}, fmt.Errorf("create data layer: %w", err)
	}
	return mutate.Addendum{
		Layer: dataLayer,
		History: v1.History{
			Author:    "sheaf",
			CreatedBy: "sheaf",
			Comment:   "experimental",
		},
	}, nil
}
