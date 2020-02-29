/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestBundle_Path(t *testing.T) {
	withBundleDir(t, func(dir string) {
		stageFile(t, "bundle.json", filepath.Join(dir, "bundle.json"))

		bundle, err := NewBundle(dir)
		require.NoError(t, err)

		actual := bundle.Path()
		require.Equal(t, dir, actual)
	})
}

func TestBundle_Config(t *testing.T) {
	withBundleDir(t, func(dir string) {
		configRaw := stageFile(t, "bundle.json", filepath.Join(dir, "bundle.json"))
		var wanted sheaf.BundleConfig
		require.NoError(t, json.Unmarshal(configRaw, &wanted))

		bundle, err := NewBundle(dir)
		require.NoError(t, err)

		actual := bundle.Config()
		require.Equal(t, wanted, actual)
	})
}

func TestBundle_Manifests(t *testing.T) {
	withBundleDir(t, func(dir string) {
		stageFile(t, "bundle.json", filepath.Join(dir, "bundle.json"))

		manifestDir := filepath.Join("testdata", "manifests")
		bundle, err := NewBundle(dir, ManifestsDirOption(manifestDir))
		require.NoError(t, err)

		actual, err := bundle.Manifests()
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
