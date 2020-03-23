/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/internal/goutil"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// BundleConfigWriterOption is an option for configuring BundleConfigWriter.
type BundleConfigWriterOption func(bcw *BundleConfigWriter)

// BundleConfigWriter writes a bundle config to the filesystem.
type BundleConfigWriter struct {
	codec    sheaf.BundleConfigCodec
	openFile func(string) (io.WriteCloser, error)
}

var _ sheaf.BundleConfigWriter = &BundleConfigWriter{}

// NewBundleConfigWriter creates an instance of BundleConfigWriter.
func NewBundleConfigWriter(options ...BundleConfigWriterOption) *BundleConfigWriter {
	bcw := BundleConfigWriter{
		codec:    NewBundleConfigCodec(),
		openFile: openFile,
	}

	for _, option := range options {
		option(&bcw)
	}

	return &bcw
}

// Write writes a bundle to the filesystem. It will use the path from the bundle argument.
func (b BundleConfigWriter) Write(bundle sheaf.Bundle, config sheaf.BundleConfig) error {
	filename := filepath.Join(bundle.Path(), sheaf.BundleConfigFilename)

	f, err := b.openFile(filename)
	if err != nil {
		return fmt.Errorf("open bundle config %q: %w", filename, err)
	}

	defer goutil.Close(f)

	if err := b.codec.Encode(f, config); err != nil {
		return fmt.Errorf("encode config: %w", err)
	}

	return nil
}

func openFile(filename string) (io.WriteCloser, error) {
	return os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
}
