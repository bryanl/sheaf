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
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/containerd/continuity/fs"

	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// ManifestServiceOption is a functional option for configuration ManifestService.
type ManifestServiceOption func(m ManifestService) ManifestService

// ManifestServiceReporter sets the reporter.
func ManifestServiceReporter(r reporter.Reporter) ManifestServiceOption {
	return func(m ManifestService) ManifestService {
		m.reporter = r
		return m
	}
}

// ManifestService is a service for interacting with manifests on a filesystem.
type ManifestService struct {
	manifestsDir string
	reporter     reporter.Reporter
}

var _ sheaf.ManifestService = &ManifestService{}

// NewManifestService creates an instance of ManifestService.
func NewManifestService(manifestsDir string, options ...ManifestServiceOption) (*ManifestService, error) {
	m := ManifestService{
		manifestsDir: manifestsDir,
		reporter:     reporter.Default,
	}

	for _, option := range options {
		m = option(m)
	}

	return &m, nil
}

// List lists manifests on the filesystem.
func (m ManifestService) List() ([]sheaf.BundleManifest, error) {
	entries, err := ioutil.ReadDir(m.manifestsDir)
	if err != nil {
		return nil, fmt.Errorf("read manifests dir %q: %w", m.manifestsDir, err)
	}

	var list []sheaf.BundleManifest

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(m.manifestsDir, entry.Name())

		data, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read manifest %q: %w", manifestPath, err)
		}

		bm := sheaf.BundleManifest{
			ID:   manifestPath,
			Data: data,
		}

		list = append(list, bm)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	return list, nil
}

// Add adds zero or more manifests to the filesystem.
func (m ManifestService) Add(manifestURIs ...string) error {
	if err := os.MkdirAll(m.manifestsDir, 0700); err != nil {
		return err
	}

	for _, manifestURI := range manifestURIs {
		m.reporter.Header(fmt.Sprintf("Adding manifest from %s", manifestURI))

		u, validURL, err := getURL(manifestURI)
		if err != nil {
			return err
		}

		if validURL {
			if err := m.addURL(*u); err != nil {
				return err
			}
			continue
		}

		manifestURI, err = filepath.Abs(manifestURI)
		if err != nil {
			return err
		}

		fi, err := os.Stat(manifestURI)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if err := m.addDir(manifestURI); err != nil {
				return err
			}
			continue
		}

		if err := m.addFile(manifestURI); err != nil {
			return err
		}
	}

	return nil
}

func (m ManifestService) addURL(manifestURL url.URL) error {
	if !strings.HasPrefix(manifestURL.Scheme, "http") {
		return fmt.Errorf("%s is an unsupported URL", manifestURL.String())
	}

	resp, err := http.Get(manifestURL.String())
	if err != nil {
		return err
	}

	defer func() {
		if cErr := resp.Body.Close(); cErr != nil {
			log.Printf("unable to close http body: %v", err)
		}
	}()

	dir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return err
	}

	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			log.Printf("unable to remove temporary directory: %v", rErr)
		}
	}()

	_, file := path.Split(manifestURL.Path)

	dest := filepath.Join(dir, file)
	f, err := os.Create(dest)
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		if cErr := f.Close(); cErr != nil {
			log.Printf("unable to close temporary file: %v", cErr)
		}
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return m.addFile(dest)

}

func (m ManifestService) addFile(manifestURI string) error {
	_, file := filepath.Split(manifestURI)

	dest := filepath.Join(m.manifestsDir, file)
	if _, err := os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			return fs.CopyFile(dest, manifestURI)
		}
		return fmt.Errorf("destination is invalid: %w", err)
	}

	return fmt.Errorf("%s exists", dest)
}

func (m ManifestService) addDir(manifestDir string) error {
	fis, err := ioutil.ReadDir(manifestDir)
	if err != nil {
		return err
	}

	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}

		manifestPath := filepath.Join(manifestDir, fi.Name())
		if err := m.addFile(manifestPath); err != nil {
			return err
		}
	}
	return nil
}

// getURL returns a url.URL if the given URI is a valid URL.
// url.Parse returns a url object with the scheme populated for
// Winddows paths, so the Host must also be checked.
func getURL(manifestURI string) (*url.URL, bool, error) {
	u, err := url.Parse(manifestURI)
	if err != nil {
		return nil, false, err
	}

	if u.Scheme != "" && u.Host != "" {
		return u, true, nil
	}
	return nil, false, nil
}
