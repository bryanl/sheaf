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

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// BundleCreator creates a function that can create bundle configs.
func BundleCreator(bundlePath string, options ...Option) func(bc sheaf.BundleConfig) error {
	opts := makeCreateBundleOptions(options...)

	return func(bc sheaf.BundleConfig) error {
		if bundlePath == "" {
			return fmt.Errorf("bundle path is blank")
		}

		if err := os.MkdirAll(bundlePath, 0700); err != nil {
			return fmt.Errorf("create bundle directory %q: %w", bundlePath, err)
		}

		bundleConfigPath := filepath.Join(bundlePath, sheaf.BundleConfigFilename)
		f, err := os.Create(bundleConfigPath)
		if err != nil {
			return fmt.Errorf("create bundle config: %w", err)
		}

		defer func() {
			if cErr := f.Close(); cErr != nil {
				log.Printf("close bundle config: %v", err)
			}
		}()

		if err := opts.bundleConfigCodec.Encode(f, bc); err != nil {
			return fmt.Errorf("encode bundle config: %w", err)
		}

		manifestsPath := filepath.Join(bundlePath, "app", "manifests")
		if err := os.MkdirAll(manifestsPath, 0700); err != nil {
			return fmt.Errorf("create manifests directory: %w", err)
		}

		return nil
	}
}
