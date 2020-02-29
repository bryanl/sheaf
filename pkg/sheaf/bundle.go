/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/images"
)

//go:generate mockgen -destination=../mocks/mock_bundle.go -package mocks github.com/bryanl/sheaf/pkg/sheaf Bundle

// ManifestGenerator generates manifests.
type ManifestGenerator interface {
	Generate(w io.Writer) error
}

// BundleFactory is a factory for creating bundles given a URI.
type BundleFactory func(uri string) (Bundle, error)

// Bundle manages bundles.
type Bundle interface {
	Codec() Codec
	Path() string
	Config() BundleConfig
	Artifacts() ArtifactsService
	Manifests() ([]BundleManifest, error)
	Images() (images.Set, error)
}

// BundleManifest describes a manifest in a bundle.
type BundleManifest struct {
	ID   string
	Data []byte
}

func loadBundleConfig(path string) (BundleConfig, string, error) {
	bundleConfig := BundleConfig{}

	// check if directory exists
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return bundleConfig, "", fmt.Errorf("bundle directory %q does not exist", path)
		}

		return bundleConfig, "", err
	}

	if !fi.IsDir() {
		return bundleConfig, "", fmt.Errorf("%q is not a directory", path)
	}

	bundleConfigFilename := filepath.Join(path, BundleConfigFilename)

	bundleConfig, err = LoadBundleConfig(bundleConfigFilename)
	if err != nil {
		return bundleConfig, "", fmt.Errorf("load bundle config: %w", err)
	}

	return bundleConfig, bundleConfigFilename, err
}
