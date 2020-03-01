/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ManifestAdderOption configures ManifestAdder.
type ManifestAdderOption func(ma ManifestAdder) ManifestAdder

// ManifestAdderForce configures ManifestAdder to override files.
func ManifestAdderForce(force bool) ManifestAdderOption {
	return func(ma ManifestAdder) ManifestAdder {
		ma.Force = force
		return ma
	}
}

// URLFetcher fetches a URL and returns a read closer, the name of the file, or a possible error.
type URLFetcher func(string) (io.ReadCloser, string, error)

// ManifestAdder adds manifests to a fs.
type ManifestAdder struct {
	// BundlePath is the path to the fs.
	BundlePath string
	// Force is set to true if files can be overwritten.
	Force bool
	// URLFetcher fetches a URL.
	URLFetcher URLFetcher
}

// NewManifestAdder creates an instance of ManifestAdder.
func NewManifestAdder(bundlePath string, options ...ManifestAdderOption) (*ManifestAdder, error) {
	if err := ensureBundlePath(bundlePath); err != nil {
		return nil, fmt.Errorf("ensure fs path %q exists: %w", bundlePath, err)
	}

	ma := ManifestAdder{
		BundlePath: bundlePath,
		URLFetcher: fetchURL,
	}

	for _, option := range options {
		ma = option(ma)
	}

	return &ma, nil
}

// Add adds manifests.
func (ma *ManifestAdder) Add(filenames []string) error {
	manifestsDir := filepath.Join(ma.BundlePath, "app", "manifests")
	if err := os.MkdirAll(manifestsDir, 0700); err != nil {
		return err
	}

	for _, filename := range filenames {
		var rc io.ReadCloser
		var base string
		var err error

		if strings.HasPrefix(filename, "http") {
			rc, base, err = ma.URLFetcher(filename)
		} else {
			rc, base, err = loadFile(filename)
		}

		if err != nil {
			return err
		}

		manifestPath := filepath.Join(manifestsDir, base)

		if !ma.Force {
			_, err = os.Stat(manifestPath)
			if !os.IsNotExist(err) {
				return fmt.Errorf("destination %q exists", manifestPath)
			}
		}

		f, err := os.OpenFile(manifestPath, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return fmt.Errorf("open destination %q: %w", manifestPath, err)
		}
		if _, err := io.Copy(f, rc); err != nil {
			return err
		}

		if err := rc.Close(); err != nil {
			return fmt.Errorf("close %q: %w", filename, err)
		}
	}

	return nil
}

// nolint:bodyclose
func fetchURL(u string) (io.ReadCloser, string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return nil, "", fmt.Errorf("http get %q: %w", u, err)
	}

	base := path.Base(u)
	return resp.Body, base, nil
}

func loadFile(filename string) (io.ReadCloser, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}

	base := filepath.Base(filename)
	return f, base, nil
}

func ensureBundlePath(bundlePath string) error {
	fi, err := os.Stat(bundlePath)
	if err == nil {
		if !fi.IsDir() {
			return fmt.Errorf("is not a directory")
		}

		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("is invalid: %w", err)
	}

	return nil
}
