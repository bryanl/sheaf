// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func initBundle(t *testing.T, h *harness, name, wd string) string {
	err := h.runSheaf(wd, defaultSheafRunSettings, "init", name)
	require.NoError(t, err, "initialize sheaf bundle")

	return filepath.Join(wd, name)
}
