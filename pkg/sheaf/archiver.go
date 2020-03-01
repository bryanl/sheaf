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
	Unarchive(src, dest string) error
	Archive(src string, w io.Writer) error
}
