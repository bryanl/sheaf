/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestManifestService_List(t *testing.T) {
	withBundleDir(t, func(dir string) {
		stageFile(t, sheaf.BundleConfigFilename, filepath.Join(dir, sheaf.BundleConfigFilename))

		manifestDir := filepath.Join("testdata", "manifests")
		m, err := NewManifestService(manifestDir)
		require.NoError(t, err)

		actual, err := m.List()
		require.NoError(t, err)

		wanted := []sheaf.BundleManifest{
			{
				ID:   filepath.Join(manifestDir, "deploy.yaml"),
				Data: slurpData(t, filepath.Join(manifestDir, "deploy.yaml")),
			},
			{
				ID:   filepath.Join(manifestDir, "service.yaml"),
				Data: slurpData(t, filepath.Join(manifestDir, "service.yaml")),
			},
		}
		require.Equal(t, wanted, actual)
	})
}

func TestManifestService_Add(t *testing.T) {
	cases := []struct {
		name         string
		manifestPath string
		setup        func(t *testing.T, bundleDir string)
		wantErr      bool
	}{
		{
			name:         "add file",
			manifestPath: filepath.Join("testdata", "manifests", "deploy.yaml"),
		},
		{
			name:         "add file (already exists)",
			manifestPath: filepath.Join("testdata", "manifests", "deploy.yaml"),
			setup: func(t *testing.T, bundleDir string) {
				dir := genManifestDir(bundleDir)
				require.NoError(t, os.MkdirAll(dir, 0700))

				f, err := os.Create(filepath.Join(dir, "deploy.yaml"))
				require.NoError(t, err)
				defer require.NoError(t, f.Close())
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withBundleDir(t, func(bundleDir string) {
				stageFile(t, sheaf.BundleConfigFilename, filepath.Join(bundleDir, sheaf.BundleConfigFilename))

				if tc.setup != nil {
					tc.setup(t, bundleDir)
				}

				manifestDir := filepath.Join(bundleDir, "app", "manifests")
				m, err := NewManifestService(manifestDir)
				require.NoError(t, err)

				err = m.Add(tc.manifestPath)
				if tc.wantErr {
					require.Error(t, err)
					return
				}

				require.NoError(t, err)
				_, err = os.Stat(filepath.Join(manifestDir, "deploy.yaml"))
				require.NoError(t, err)
			})

		})
	}
}
