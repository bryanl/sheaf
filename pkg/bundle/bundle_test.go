/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestBundle_Path(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(dir))
	}()

	stageFile(t, "bundle.json", filepath.Join(dir, "bundle.json"))

	bundle, err := NewBundle(dir)
	require.NoError(t, err)

	actual := bundle.Path()
	require.Equal(t, dir, actual)
}

func TestBundle_Config(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(dir))
	}()

	configRaw := stageFile(t, "bundle.json", filepath.Join(dir, "bundle.json"))
	var wanted sheaf.BundleConfig
	require.NoError(t, json.Unmarshal(configRaw, &wanted))

	bundle, err := NewBundle(dir)
	require.NoError(t, err)

	actual := bundle.Config()
	require.Equal(t, wanted, actual)

}
