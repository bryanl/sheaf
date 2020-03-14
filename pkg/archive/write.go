/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/fs"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// Write writes an archive's configuration to a registry.
func Write(archivePath, dest string, forceInsecure bool) error {
	w := newWriter()
	return w.Write(archivePath, dest, forceInsecure)
}

type writeOption func(w writer) writer

func writerOptionArchiver(a sheaf.Archiver) writeOption {
	return func(w writer) writer {
		w.archiver = a
		return w
	}
}

type fsWriter func(bundlePath, dest string, forceInsecure bool) error

func writerOptionFSWriter(fw fsWriter) writeOption {
	return func(w writer) writer {
		w.fsWriter = fw
		return w
	}
}

type writer struct {
	archiver sheaf.Archiver
	reporter reporter.Reporter
	fsWriter func(bundlePath, dest string, forceInsecure bool) error
}

func newWriter(options ...writeOption) *writer {
	w := writer{
		archiver: archiver.Default,
		reporter: reporter.Default,
		fsWriter: fs.Write,
	}

	for _, option := range options {
		w = option(w)
	}

	return &w
}

func (w *writer) Write(archivePath, dest string, forceInsecure bool) error {
	w.reporter.Header("Staging archive to temporary directory")

	dir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}

	defer func() {
		if rErr := os.RemoveAll(dir); rErr != nil {
			log.Printf("remove temporary directory: %v", err)
		}
	}()

	if err := w.archiver.Unarchive(archivePath, dir); err != nil {
		return fmt.Errorf("unarchive: %w", err)
	}

	if err := os.RemoveAll(filepath.Join(dir, "artifacts")); err != nil {
		return fmt.Errorf("unable to clean artifacts directory: %w", err)
	}

	return w.fsWriter(dir, dest, forceInsecure)
}
