/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import "io"

//go:generate mockgen -destination=../mocks/mock_archiver.go -package mocks github.com/bryanl/sheaf/pkg/sheaf Archiver

// Archiver manages archives.
type Archiver interface {
	// Archive archives a source to a writer.
	Archive(src string, w io.Writer) error
	// Unarchive unarchives a reader to a destination.
	Unarchive(r io.Reader, dest string) error
	// UnarchivePath unarchives a source to a destination.
	UnarchivePath(src string, dest string) error
}
