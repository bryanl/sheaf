/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archiver

import (
	"io"
	"log"
	"os"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

var (
	// Default is the default archiver. It assumes the archive is tar.gz.
	Default = NewArchiver()
)

// Archiver manages tar.gz archives.
type Archiver struct{}

var _ sheaf.Archiver = &Archiver{}

// NewArchiver creates an instance of Archiver.
func NewArchiver() *Archiver {
	a := Archiver{}

	return &a
}

// Archive creates a gzipped tarball.
func (a Archiver) Archive(src string, w io.Writer) error {
	return targz(src, w)
}

// Unarchive unarchives a tar.gz file from src to dest.
func (a Archiver) Unarchive(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("close archive: %v", cErr)
		}
	}()

	return untargz(f, dest)
}
