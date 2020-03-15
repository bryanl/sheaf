// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/pathmapping"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

func Test_sheaf_archive_pack(t *testing.T) {
	withWorkingDirectory(t, func(options wdOptions) {
		b := sheafInit(t, testHarness, "integration", options.dir)

		require.NoError(t, b.harness.runSheaf(b.dir, defaultSheafRunSettings, "manifest", "add",
			"-f", testdata(t, "archive", "pack")))

		require.NoError(t, b.harness.runSheaf(options.dir, defaultSheafRunSettings, "archive", "pack",
			"--bundle-path", b.dir))

		var buf bytes.Buffer
		settings := genSheafRunSettings()
		settings.Stdout = &buf

		require.NoError(t, b.harness.runSheaf(options.dir, settings, "archive", "list-images",
			"integration-0.1.0.tgz"))

		var list []sheaf.BundleImage
		require.NoError(t, json.Unmarshal(buf.Bytes(), &list))

		require.Len(t, list, 1)
		require.Equal(t, "docker.io/bryanl/slim-hello-world:v1", list[0].Name)
	})
}

func Test_sheaf_archive_push(t *testing.T) {
	withWorkingDirectory(t, func(options wdOptions) {
		b := sheafInit(t, testHarness, "integration", options.dir)

		settings := defaultSheafRunSettings

		require.NoError(t, b.harness.runSheaf(b.dir, settings, "manifest", "add",
			"-f", testdata(t, "archive", "push")))

		require.NoError(t, b.harness.runSheaf(options.dir, defaultSheafRunSettings, "archive", "pack",
			"--bundle-path", b.dir))

		archivePath := filepath.Join(options.dir, b.archiveName())
		ref := genRegistryPath(options)
		pushArgs := append([]string{"archive", "push", archivePath, ref, "--insecure-registry"})
		require.NoError(t, b.harness.runSheaf(b.dir, settings, pushArgs...))

		dir, err := ioutil.TempDir("", "sheaf-test")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.RemoveAll(dir))
		}()

		dest := filepath.Join(dir, "dest")

		pullArgs := append([]string{"config", "pull", ref, dest, "--insecure-registry"})
		require.NoError(t, b.harness.runSheaf(b.dir, settings, pullArgs...))

		checkBundleEquals(t, b, dest)
	})
}

func Test_sheaf_archive_relocate(t *testing.T) {
	withWorkingDirectory(t, func(options wdOptions) {
		b := sheafInit(t, testHarness, "integration", options.dir)

		settings := defaultSheafRunSettings

		require.NoError(t, b.harness.runSheaf(b.dir, settings, "manifest", "add",
			"-f", testdata(t, "archive", "relocate")))

		require.NoError(t, b.harness.runSheaf(options.dir, settings, "archive", "pack",
			"--bundle-path", b.dir))

		archivePath := filepath.Join(options.dir, b.archiveName())
		ref := genRegistryRoot(options)
		require.NoError(t, b.harness.runSheaf(options.dir, settings, "archive", "relocate",
			archivePath, ref, "--insecure-registry"))

		originalName, err := image.NewName("docker.io/bryanl/slim-hello-world")
		require.NoError(t, err, "unable to generate name")

		expectedRepo, err := pathmapping.FlattenRepoPathPreserveTagDigest(ref, originalName)
		require.NoError(t, err)

		repo, err := name.NewRepository(expectedRepo.String(), name.Insecure)
		require.NoError(t, err)

		tags, err := remote.List(repo)
		require.NoError(t, err)

		require.Contains(t, tags, "v1")
	})
}
