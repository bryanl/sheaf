/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fsutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsDirectory(t *testing.T) {
	dir, err := ioutil.TempDir("", "sheaf-test")
	require.NoError(t, err)

	file := filepath.Join(dir, "file.txt")
	err = ioutil.WriteFile(file, []byte{}, 0644)
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(dir))
	}()

	cases := []struct {
		name     string
		in       string
		expected bool
	}{
		{
			name:     "a directory",
			in:       dir,
			expected: true,
		},
		{
			name:     "a file",
			in:       file,
			expected: false,
		},
		{
			name:     "invalid",
			in:       filepath.Join(dir, "invalid"),
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.expected, IsDirectory(tc.in))
		})
	}
}
