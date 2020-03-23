/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archiver

import (
	"fmt"
	"io"
	"os"

	"github.com/bryanl/sheaf/internal/goutil"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// Archiver is a targz archiver.
type Archiver struct{}

var _ sheaf.Archiver = &Archiver{}

// New creates an instance of Archiver.
func New() *Archiver {
	a := Archiver{}

	return &a
}

// Unarchive unarchives a reader to a directory
func (a Archiver) Unarchive(r io.Reader, dest string) error {
	return untargz(r, dest)
}

// Archive archives a directory to a writer.
func (a Archiver) Archive(src string, w io.Writer) error {
	return targz(src, w)
}

// UnarchivePath unarchives a targz file to a directory.
func (a Archiver) UnarchivePath(src string, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to open archive %q: %q", src, err)
	}

	defer goutil.Close(f)

	if err := a.Unarchive(f, dest); err != nil {
		return err
	}

	return nil
}
