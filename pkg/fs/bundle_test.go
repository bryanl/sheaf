/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestBundle_Path(t *testing.T) {
	testutil.WithBundleDir(t, func(dir string) {
		testutil.StageFile(t, sheaf.BundleConfigFilename, filepath.Join(dir, sheaf.BundleConfigFilename))

		bundle, err := NewBundle(dir)
		require.NoError(t, err)

		actual := bundle.Path()
		require.Equal(t, dir, actual)
	})
}

func TestBundle_Config(t *testing.T) {
	testutil.WithBundleDir(t, func(dir string) {
		configRaw := testutil.StageFile(t, sheaf.BundleConfigFilename, filepath.Join(dir, sheaf.BundleConfigFilename))
		var wanted sheaf.BundleConfig
		require.NoError(t, json.Unmarshal(configRaw, &wanted))

		bundle, err := NewBundle(dir)
		require.NoError(t, err)

		actual := bundle.Config()
		require.Equal(t, wanted, actual)
	})
}
