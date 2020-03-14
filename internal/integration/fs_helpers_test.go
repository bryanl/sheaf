// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pivotal/go-ape/pkg/filecopy"
	"github.com/stretchr/testify/require"
)

func withWorkingDirectory(t *testing.T, fn func(dir string)) {
	workingDirectory, err := ioutil.TempDir("", "sheaf-test")
	require.NoError(t, err)

	t.Cleanup(func() {
		if workingDirectory != "" {
			require.NoError(t, os.RemoveAll(workingDirectory))
		}
	})

	fn(workingDirectory)
}

func stageFile(t *testing.T, srcPath, destPath string) {
	err := filecopy.Copy(destPath, srcPath)
	require.NoError(t, err)
}

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

func checkFileMatches(t *testing.T, file string, want []byte) {
	actual, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	require.Equal(t, string(want), string(actual))
}

func checkFileEquals(t *testing.T, file1, file2 string) {
	b := readFile(t, file1)
	checkFileMatches(t, file2, b)
}

func readFile(t *testing.T, file string) []byte {
	data, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	return data
}

func readJSONFile(t *testing.T, file string, v interface{}) {
	data := readFile(t, file)
	err := json.Unmarshal(data, v)
	require.NoError(t, err)
}

func writeJSONFile(t *testing.T, file string, v interface{}) {
	bytes, err := json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)

	err = ioutil.WriteFile(file, bytes, 0644)
	require.NoError(t, err)
}

func testdata(t *testing.T, parts ...string) string {
	f, err := filepath.Abs(filepath.Join(append([]string{"testdata"}, parts...)...))
	require.NoError(t, err)

	return f
}
