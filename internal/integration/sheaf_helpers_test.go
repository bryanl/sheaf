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

	"github.com/bryanl/sheaf/pkg/sheaf"
)

func sheafInit(t *testing.T, h *harness, name, wd string) *bundle {
	err := h.runSheaf(wd, defaultSheafRunSettings, "init", name)
	require.NoError(t, err, "initialize sheaf bundle")

	b := bundle{
		dir:     filepath.Join(wd, name),
		harness: h,
	}

	return &b
}

type bundle struct {
	dir     string
	harness *harness
}

func (b bundle) readConfig(t *testing.T) sheaf.BundleConfig {
	var config sheaf.BundleConfig
	readJSONFile(t, b.configFile(), &config)
	return config
}

func (b bundle) updateConfig(t *testing.T, fn func(config *sheaf.BundleConfig)) {
	var config sheaf.BundleConfig
	readJSONFile(t, b.configFile(), &config)

	fn(&config)
	writeJSONFile(t, b.configFile(), config)
}

func (b bundle) configFile() string {
	return filepath.Join(b.dir, "bundle.json")
}
