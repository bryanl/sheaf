/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/codec"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// Option is a functional option for configuring Bundle.
type Option func(b Bundle) Bundle

// CodecOption sets the codec for the bundle.
func CodecOption(c sheaf.Codec) Option {
	return func(b Bundle) Bundle {
		b.codec = c
		return b
	}
}

// Bundle is a bundle that lives on a filesystem.
type Bundle struct {
	rootPath string
	config   sheaf.BundleConfig
	codec    sheaf.Codec
}

var _ sheaf.BundleService = &Bundle{}

// NewBundle creates an instance of Bundle. `rootPath` points to root directory
// of the bundle on the filesystem.
func NewBundle(rootPath string, options ...Option) (*Bundle, error) {
	config, err := loadBundleConfig(rootPath)
	if err != nil {
		return nil, fmt.Errorf("load bundle config: %w", err)
	}

	b := Bundle{
		rootPath: rootPath,
		config:   config,
	}

	for _, option := range options {
		b = option(b)
	}

	if b.codec == nil {
		b.codec = codec.Default
	}

	return &b, nil
}

// Artifacts returns an artifacts service for the bundle.
func (b *Bundle) Artifacts() sheaf.ArtifactsService {
	return NewArtifactsService(b)
}

// Path returns the root path of the bundle.
func (b *Bundle) Path() string {
	return b.rootPath
}

// Config returns the configuration for the bundle.
func (b *Bundle) Config() sheaf.BundleConfig {
	return b.config
}

func loadBundleConfig(path string) (sheaf.BundleConfig, error) {
	bundleConfig := sheaf.BundleConfig{}

	// check if directory exists
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return bundleConfig, fmt.Errorf("bundle directory %q does not exist", path)
		}

		return bundleConfig, err
	}

	if !fi.IsDir() {
		return bundleConfig, fmt.Errorf("%q is not a directory", path)
	}

	bundleConfigFilename := filepath.Join(path, sheaf.BundleConfigFilename)

	bundleConfig, err = sheaf.LoadBundleConfig(bundleConfigFilename)
	if err != nil {
		return bundleConfig, fmt.Errorf("load bundle config: %w", err)
	}

	return bundleConfig, err
}

// Codec is the codec for the bundle.
func (b *Bundle) Codec() sheaf.Codec {
	return b.codec
}
