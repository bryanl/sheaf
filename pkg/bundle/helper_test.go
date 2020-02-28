/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle

import (
	"io/ioutil"
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
