/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fsutil

import "os"

// IsDirectory returns true if the path is a directory
func IsDirectory(p string) bool {
	fi, err := os.Stat(p)
	if err != nil {
		return false
	}

	return fi.IsDir()
}
