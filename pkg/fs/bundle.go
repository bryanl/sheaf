/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bryanl/sheaf/pkg/codec"
	"github.com/bryanl/sheaf/pkg/images"
	"github.com/bryanl/sheaf/pkg/manifest"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// DefaultBundleConfigWriter is the default bundle config writer.
func DefaultBundleConfigWriter(bundle sheaf.Bundle, config sheaf.BundleConfig) error {
	jbc, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	filename := filepath.Join(bundle.Path(), sheaf.BundleConfigFilename)

	return ioutil.WriteFile(filename, jbc, 0644)
}

// DefaultBundleFactory is the default fs factory.
func DefaultBundleFactory(uri string) (sheaf.Bundle, error) {
	return NewBundle(uri)
}

// Option is a functional option for configuring Bundle.
type Option func(b Bundle) Bundle

// CodecOption sets the codec for the fs.
func CodecOption(c sheaf.Codec) Option {
	return func(b Bundle) Bundle {
		b.codec = c
		return b
	}
}

// ManifestsDirOption sets the location to the fs's manifest.
func ManifestsDirOption(p string) Option {
	return func(b Bundle) Bundle {
		b.manifestsDir = p
		return b
	}
}

// Bundle is a fs that lives on a filesystem.
type Bundle struct {
	rootPath     string
	config       sheaf.BundleConfig
	codec        sheaf.Codec
	manifestsDir string
	reporter     reporter.Reporter
}

var _ sheaf.Bundle = &Bundle{}

// NewBundle creates an instance of Bundle. `rootPath` points to root directory
// of the fs on the filesystem.
func NewBundle(rootPath string, options ...Option) (*Bundle, error) {
	rootPath, err := locateRootDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("locate bundle root directory")
	}

	config, err := loadBundleConfig(rootPath)
	if err != nil {
		return nil, fmt.Errorf("load fs config: %w", err)
	}

	b := Bundle{
		rootPath: rootPath,
		config:   config,
		reporter: reporter.Default,
	}

	for _, option := range options {
		b = option(b)
	}

	if b.codec == nil {
		b.codec = codec.Default
	}

	if b.manifestsDir == "" {
		b.manifestsDir = filepath.Join(b.rootPath, "app", "manifests")
	}

	return &b, nil
}

// Artifacts returns an artifacts service for the fs.
func (b *Bundle) Artifacts() sheaf.ArtifactsService {
	return NewArtifactsService(b)
}

// Path returns the root path of the fs.
func (b *Bundle) Path() string {
	return b.rootPath
}

// Config returns the configuration for the fs.
func (b *Bundle) Config() sheaf.BundleConfig {
	return b.config
}

func loadBundleConfig(path string) (sheaf.BundleConfig, error) {
	bundleConfig := sheaf.BundleConfig{}

	path, err := locateRootDir(path)
	if err != nil {
		return bundleConfig, err
	}

	// check if directory exists
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return bundleConfig, fmt.Errorf("fs directory %q does not exist", path)
		}

		return bundleConfig, err
	}

	if !fi.IsDir() {
		return bundleConfig, fmt.Errorf("%q is not a directory", path)
	}

	bundleConfigFilename := filepath.Join(path, sheaf.BundleConfigFilename)

	data, err := ioutil.ReadFile(bundleConfigFilename)
	if err != nil {
		return sheaf.BundleConfig{}, fmt.Errorf("read %q: %w", bundleConfigFilename, err)
	}

	var bc sheaf.BundleConfig
	if err := json.Unmarshal(data, &bc); err != nil {
		return sheaf.BundleConfig{}, err
	}

	return bc, nil
}

// Codec is the codec for the fs.
func (b *Bundle) Codec() sheaf.Codec {
	return b.codec
}

// Manifests returns the manifest service for the bundle.
func (b *Bundle) Manifests() (sheaf.ManifestService, error) {
	manifestsDir, err := locateManifestsDir(b.rootPath)
	if err != nil {
		return nil, fmt.Errorf("locate manifest directory: %w", err)
	}

	return NewManifestService(manifestsDir)
}

func locateRootDir(in string) (string, error) {
	in, err := filepath.Abs(in)
	if err != nil {
		return "", err
	}

	configPath := filepath.Join(in, sheaf.BundleConfigFilename)
	if _, err := os.Stat(configPath); err != nil {
		if !os.IsNotExist(err) {
			return "", err
		}

		in = filepath.Clean(in)
		if strings.HasSuffix(in, string(filepath.Separator)) {
			return "", fmt.Errorf("bundle config not found")
		}

		dir := filepath.Dir(in)
		return locateRootDir(dir)
	}

	return in, nil
}

// locate a manifest directory given a path.
// TODO: ensure this works on windows
func locateManifestsDir(in string) (string, error) {
	rootDir, err := locateRootDir(in)
	if err != nil {
		return "", fmt.Errorf("locate bundle root directory: %w", err)
	}

	return genManifestDir(rootDir), nil
}

func genManifestDir(rootPath string) string {
	return filepath.Join(rootPath, "app", "manifests")
}

// Images returns images in the fs.
func (b *Bundle) Images() (images.Set, error) {
	seen := images.Empty

	config := b.Config()
	bundleImages := *config.Images
	printImageTree(sheaf.BundleConfigFilename, bundleImages.Strings(), b.reporter)

	seen = seen.Union(bundleImages)

	m, err := b.Manifests()
	if err != nil {
		return images.Empty, err
	}

	bundleManifests, err := m.List()
	if err != nil {
		return images.Empty, err
	}

	for _, bundleManifest := range bundleManifests {

		list, err := manifest.ContainerImages(bundleManifest.ID, config.UserDefinedImages)
		if err != nil {
			return images.Empty, fmt.Errorf("find container images for %s: %w", bundleManifest, err)
		}

		names := list.Strings()
		if len(names) < 1 {
			continue
		}

		p := strings.TrimPrefix(bundleManifest.ID, b.manifestsDir+"/")
		printImageTree(p, names, b.reporter)

		seen = seen.Union(list)
	}

	return seen, nil
}

func printImageTree(source string, imageNames []string, r reporter.Reporter) {
	if len(imageNames) == 0 {
		return
	}
	r.Report(source)
	for i, name := range imageNames {
		prefix := treeItem
		if i == len(imageNames)-1 {
			prefix = treeItemLast
		}

		r.Reportf("%s %s", prefix, name)
	}
}
