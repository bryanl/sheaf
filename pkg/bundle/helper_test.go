/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func stageFile(t *testing.T, name, dest string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	require.NoError(t, ioutil.WriteFile(dest, data, 0600))

	return data
}

func withBundleDir(t *testing.T, fn func(dir string)) {
	dir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(dir))
	}()

	fn(dir)
}

func slurpData(t *testing.T, p string) []byte {
	data, err := ioutil.ReadFile(p)
	require.NoError(t, err)
	return data
}
