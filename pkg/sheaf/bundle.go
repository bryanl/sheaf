/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cnabio/duffle/pkg/imagestore/ocilayout"
	"github.com/deislabs/duffle/pkg/imagestore"
	dcopy "github.com/otiai10/copy"
)

// Bundle represents a bundle
type Bundle struct {
	// Path is the path to the bundle directory.
	Path string
	// Config is the BundleConfig for the bundle.
	Config BundleConfig
	// Store is the image store
	Store imagestore.Store

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
	if err := dcopy.Copy(path, root); err != nil {
		return nil, fmt.Errorf("stage bundle: %w", err)
	}

	store, err := ocilayout.Create(setStoreLocation(root))
	if err != nil {
		return nil, fmt.Errorf("create image store: %w", err)
	}

	bundle := &Bundle{
		Path:   root,
		Config: bundleConfig,
		Store:  store,
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

func setStoreLocation(archiveDir string) imagestore.Option {
	return func(parameters imagestore.Parameters) imagestore.Parameters {
		parameters.ArchiveDir = archiveDir
		return parameters
	}
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
func (b *Bundle) Images() ([]string, error) {
	seen := make(map[string]bool)

	manifestPaths, err := b.Manifests()
	if err != nil {
		return nil, err
	}

	for _, manifestPath := range manifestPaths {
		images, err := ContainerImages(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("find container images for %q: %w", manifestPath, err)
		}

		fmt.Printf("Images in %s: [%s]\n",
			filepath.Base(manifestPath), strings.Join(images, ","))
		for i := range images {
			seen[images[i]] = true
		}
	}

	var list []string
	for k := range seen {
		list = append(list, k)
	}

	return union(list, b.Config.Images), nil
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
