/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pivotal/go-ape/pkg/filecopy"
	"github.com/stretchr/testify/require"
)

type wdOptions struct {
	dir      string
	registry string
}

func withWorkingDirectory(t *testing.T, fn func(options wdOptions)) {
	workingDirectory, err := ioutil.TempDir("", "sheaf-test")
	require.NoError(t, err)

	t.Cleanup(func() {
		if workingDirectory != "" {
			require.NoError(t, os.RemoveAll(workingDirectory))
		}
	})

	options := wdOptions{
		dir: workingDirectory,
	}

	fn(options)
}

func withWorkingDirectoryAndMaybeRegistry(t *testing.T, fn func(options wdOptions)) {
	// If no registry is available, skip the test.
	if os.Getenv("REGISTRY_UNAVAILABLE") != "" {
		return
	}
	workingDirectory, err := ioutil.TempDir("", "sheaf-test")
	require.NoError(t, err)

	t.Cleanup(func() {
		if workingDirectory != "" {
			require.NoError(t, os.RemoveAll(workingDirectory))
		}
	})

	registry := os.Getenv("REGISTRY")

	if registry == "" {
		r := newRegistry()
		r.Start(t)
		defer r.Stop(t)

		registry = r.Port(t)
	}

	options := wdOptions{
		dir:      workingDirectory,
		registry: registry,
	}

	fn(options)
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

	require.Equal(t, string(want), string(actual), "%s did not match wanted value", file)
}

func checkFileEquals(t *testing.T, file1, file2 string) {
	b := readFile(t, file1)
	checkFileMatches(t, file2, b)
}

func checkBundleEquals(t *testing.T, b *bundle, dest string) {
	destConfig := filepath.Join(dest, "bundle.json")
	checkFileEquals(t, b.configFile(), destConfig)

	fis, err := ioutil.ReadDir(b.pathJoin("app", "manifests"))
	require.NoError(t, err)

	for _, fi := range fis {
		cur := filepath.Join(dest, "app", "manifests", fi.Name())
		checkFileEquals(t,
			b.pathJoin("app", "manifests", fi.Name()),
			cur)
	}
}

func readFile(t *testing.T, file string) []byte {
	data, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	return data
}

func testdata(t *testing.T, parts ...string) string {
	f, err := filepath.Abs(filepath.Join(append([]string{"testdata"}, parts...)...))
	require.NoError(t, err)

	return f
}
