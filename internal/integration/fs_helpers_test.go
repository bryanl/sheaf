// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func checkDirExists(t *testing.T, dir string) {
	fi, err := os.Stat(dir)
	require.NoError(t, err, "path %s does not exist", dir)
	require.True(t, fi.IsDir(), "path %s is not a directory", dir)
}

func checkFileExists(t *testing.T, file string) {
	fi, err := os.Stat(file)
	require.NoError(t, err, "path %s does not exist", file)
	require.False(t, fi.IsDir(), "path %s is not a file", file)
}
