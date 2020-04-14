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
	"strings"

	"github.com/pivotal/go-ape/pkg/filecopy"
	"github.com/pivotal/image-relocation/pkg/images"

	"github.com/bryanl/sheaf/pkg/codec"
	"github.com/bryanl/sheaf/pkg/manifest"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// BundleOption is a functional option for configuring Bundle.
type BundleOption func(b Bundle) Bundle

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
func NewBundle(bundleDir string, options ...BundleOption) (*Bundle, error) {
	rootPath, err := locateRootDir(bundleDir)
	if err != nil {
		return nil, fmt.Errorf("locate bundle root directory for %s: %w", bundleDir, err)
	}

	config, err := LoadBundleConfig(rootPath)
	if err != nil {
		return nil, fmt.Errorf("load bundle config: %w", err)
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

// LoadBundleConfig loads a bundle config from a path on the filesystem.
func LoadBundleConfig(path string) (sheaf.BundleConfig, error) {
	path, err := locateRootDir(path)
	if err != nil {
		return nil, err
	}

	// check if directory exists
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("fs directory %q does not exist", path)
		}

		return nil, err
	}

	if !fi.IsDir() {
		return nil, fmt.Errorf("%q is not a directory", path)
	}

	bundleConfigFilename := filepath.Join(path, sheaf.BundleConfigFilename)

	f, err := os.Open(bundleConfigFilename)
	if err != nil {
		return nil, err
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("close bundle config: %v", err)
		}
	}()

	bcc := BundleConfigCodec{}
	bc, err := bcc.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode bundle: %w", err)
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
	config := b.Config()

	m, err := b.Manifests()
	if err != nil {
		return images.Empty, err
	}
	bundleManifests, err := m.List()
	if err != nil {
		return images.Empty, err
	}

	seen := images.Empty
	for _, bundleManifest := range bundleManifests {

		list, err := manifest.ContainerImages(bundleManifest.ID, config.GetUserDefinedImages())
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

// Copy copies the bundle to a new path and returns a new bundle.
func (b *Bundle) Copy(dest string) (sheaf.Bundle, error) {
	if err := filecopy.Copy(
		filepath.Join(dest, sheaf.BundleConfigFilename),
		filepath.Join(b.Path(), sheaf.BundleConfigFilename)); err != nil {
		return nil, fmt.Errorf("copy bundle config to %s: %w", dest, err)
	}

	if err := filecopy.Copy(
		filepath.Join(dest, "app", "manifests"),
		filepath.Join(b.Path(), "app", "manifests")); err != nil {
		return nil, fmt.Errorf("copy manifests to %s: %w", dest, err)
	}

	nb, err := NewBundle(dest,
		func(x Bundle) Bundle {
			x.reporter = b.reporter
			x.codec = b.codec
			return x
		},
	)

	if err != nil {
		return nil, fmt.Errorf("create new bundle: %w", err)
	}

	return nb, nil
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
