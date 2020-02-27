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
	"net/http"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/images"
	"github.com/pivotal/go-ape/pkg/filecopy"
	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
)

// Bundle represents a bundle
type Bundle struct {
	// Path is the path to the bundle directory.
	Path string
	// Config is the BundleConfig for the bundle.
	Config BundleConfig
	// Layout is the OCI image layout for the bundle.
	Layout registry.Layout

	// tmpDir for temporary things.
	tmpDir string
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

// OpenBundle loads a bundle. Call Bundle.Close() to ensure workspace is cleaned up.
func OpenBundle(path string) (*Bundle, error) {
	bundleConfig, _, err := loadBundleConfig(path)
	if err != nil {
		return nil, err
	}

	tmpDir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return nil, fmt.Errorf("create temp directory: %w", err)
	}

	root := filepath.Join(tmpDir, filepath.Base(path))
	if err := filecopy.Copy(root, path); err != nil {
		return nil, fmt.Errorf("copy bundle: %w", err)
	}

	layout, err := ggcr.NewRegistryClient(ggcr.WithTransport(http.DefaultTransport)).
		NewLayout(filepath.Join(root, "artifacts", "layout"))
	if err != nil {
		return nil, fmt.Errorf("create OCI image layout: %w", err)
	}

	bundle := &Bundle{
		Path:   root,
		Config: bundleConfig,
		Layout: layout,
		tmpDir: tmpDir,
	}

	return bundle, nil
}

// ImportBundle imports a bundle from an archive. It unpacks the bundle to a temporary
// directory.
func ImportBundle(archivePath, unpackDir string) (*Bundle, error) {
	source, err := os.Open(archivePath)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer func() {
		if cErr := source.Close(); cErr != nil {
			log.Printf("unable to close %s: %v", archivePath, err)
		}
	}()

	if err := Unarchive(source, unpackDir); err != nil {
		return nil, fmt.Errorf("unpack bundle: %w", err)
	}

	return OpenBundle(unpackDir)
}

// Manifests returns paths to manifests contained in the bundle.
// It assumes manifests live in `app/manifests`.
func (b *Bundle) Manifests() ([]string, error) {
	manifestsPath := filepath.Join(b.Path, "app", "manifests")
	entries, err := ioutil.ReadDir(manifestsPath)
	if err != nil {
		return nil, fmt.Errorf("read manifests dir %q: %w", manifestsPath, err)
	}

	var list []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(manifestsPath, entry.Name())
		list = append(list, manifestPath)
	}

	return list, nil
}

// Images returns images declared in the bundle config together with any present in manifests in the bundle.
// Images are found in manifests by searching for pod specs and iterating over the containers.
func (b *Bundle) Images() (images.Set, error) {
	seen := images.Empty

	manifestPaths, err := b.Manifests()
	if err != nil {
		return images.Empty, err
	}

	for _, manifestPath := range manifestPaths {
		imgs, err := ContainerImages(manifestPath)
		if err != nil {
			return images.Empty, fmt.Errorf("find container images for %q: %w", manifestPath, err)
		}

		printImageList(filepath.Base(manifestPath), imgs)

		seen = seen.Union(imgs)
	}

	printImageList(BundleConfigFilename, b.Config.Images)

	return seen.Union(b.Config.Images), nil
}

func printImageList(source string, imgs images.Set) {
	fmt.Printf("Images in %s: %v\n", source, imgs)
}

// Bundle writes archive to disk.
func (b *Bundle) Write() error {
	outputFile := b.Config.Filename(".")
	fmt.Println("Creating archive: ", outputFile)
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("unable to write %s: %v", outputFile, err)
		}
	}()

	if err := Archive(b.Path, f); err != nil {
		return fmt.Errorf("create archive: %w", err)
	}

	return nil
}

// Close closes the bundle and cleans up temporary files.
func (b *Bundle) Close() error {
	if err := os.RemoveAll(b.tmpDir); err != nil {
		return fmt.Errorf("remove temporary directory")
	}

	return nil
}
