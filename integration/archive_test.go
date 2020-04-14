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

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/pathmapping"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func Test_sheaf_archive_pack(t *testing.T) {
	withWorkingDirectory(t, func(options wdOptions) {
		b := sheafInit(t, testHarness, "integration", options.dir)

		_, err := b.harness.runSheaf(b.dir, "manifest", "add",
			"-f", testdata(t, "archive", "pack"))
		require.NoError(t, err)

		_, err = b.harness.runSheaf(options.dir, "archive", "pack",
			"--bundle-path", b.dir)
		require.NoError(t, err)

		listOutput, err := b.harness.runSheaf(options.dir, "archive", "list-images",
			"--archive", "integration-0.1.0.tgz")
		require.NoError(t, err)

		var list []sheaf.BundleImage
		require.NoError(t, json.Unmarshal(listOutput.Stdout.Bytes(), &list))

		require.Len(t, list, 1)
		require.Equal(t, "docker.io/bryanl/slim-hello-world:v1", list[0].Name)
	})
}

func Test_sheaf_archive_push(t *testing.T) {
	withWorkingDirectoryAndMaybeRegistry(t, func(options wdOptions) {
		b := sheafInit(t, testHarness, "integration", options.dir)

		_, err := b.harness.runSheaf(b.dir, "manifest", "add",
			"-f", testdata(t, "archive", "push"))
		require.NoError(t, err)

		_, err = b.harness.runSheaf(options.dir, "archive", "pack",
			"--bundle-path", b.dir)
		require.NoError(t, err)

		archivePath := filepath.Join(options.dir, b.archiveName())
		ref := genRegistryPath(options)
		pushArgs := []string{"archive", "push", "--archive", archivePath, "--ref", ref, "--insecure-registry"}
		_, err = b.harness.runSheaf(b.dir, pushArgs...)
		require.NoError(t, err)

		dir, err := ioutil.TempDir("", "sheaf-test")
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.RemoveAll(dir))
		}()

		dest := filepath.Join(dir, "dest")

		pullArgs := []string{"config", "pull", "--ref", ref, "--dest", dest, "--insecure-registry"}
		_, err = b.harness.runSheaf(b.dir, pullArgs...)
		require.NoError(t, err)

		checkBundleEquals(t, b, dest)
	})
}

func Test_sheaf_archive_relocate(t *testing.T) {
	withWorkingDirectoryAndMaybeRegistry(t, func(options wdOptions) {
		b := sheafInit(t, testHarness, "integration", options.dir)

		_, err := b.harness.runSheaf(b.dir, "manifest", "add",
			"-f", testdata(t, "archive", "relocate"))
		require.NoError(t, err)

		_, err = b.harness.runSheaf(options.dir, "archive", "pack",
			"--bundle-path", b.dir)
		require.NoError(t, err)

		archivePath := filepath.Join(options.dir, b.archiveName())
		ref := genRegistryRoot(options)
		_, err = b.harness.runSheaf(options.dir, "archive", "relocate",
			"--archive", archivePath, "--prefix", ref, "--insecure-registry")
		require.NoError(t, err)

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

func Test_sheaf_archive_show_manifests(t *testing.T) {
	td := func(parts ...string) string {
		return testdata(t, append([]string{"archive", "show-manifests"}, parts...)...)
	}

	cases := []struct {
		name      string
		manifests []string
		args      []string
		wanted    []byte
	}{
		{
			name: "show single manifest",
			manifests: []string{
				td("workload1.yaml"),
			},
			wanted: readFile(t, td("single.yaml")),
		},
		{
			name: "show multiple manifests",
			manifests: []string{
				td("workload2.yaml"),
				td("workload1.yaml"),
			},
			wanted: readFile(t, td("multiple.yaml")),
		},
		{
			name: "show manifests with prefix",
			manifests: []string{
				td("workload1.yaml"),
			},
			args: []string{
				"--prefix", "example.com/registry",
			},
			wanted: readFile(t, td("prefix.yaml")),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withWorkingDirectory(t, func(options wdOptions) {
				b := sheafInit(t, testHarness, "integration", options.dir)

				for _, manifest := range tc.manifests {
					_, err := b.harness.runSheaf(b.dir, "manifest", "add",
						"-f", manifest)
					require.NoError(t, err, "adding manifest %s", manifest)
				}

				_, err := b.harness.runSheaf(options.dir, "archive", "pack",
					"--bundle-path", b.dir)
				require.NoError(t, err)

				archivePath := filepath.Join(options.dir, b.archiveName())
				args := append([]string{"archive", "show-manifests", "--archive", archivePath}, tc.args...)

				output, err := b.harness.runSheaf(b.dir, args...)
				require.NoError(t, err)

				require.Equal(t, string(testutil.NormalizeNewlines(tc.wanted)), string(testutil.NormalizeNewlines(output.Stdout.Bytes())))
			})

		})
	}

}
