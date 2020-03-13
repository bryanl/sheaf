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

	"github.com/bryanl/sheaf/pkg/images"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func Test_sheaf_config_add_image(t *testing.T) {
	cases := []struct {
		name    string
		initial []string
		images  []string
		wanted  []string
	}{
		{
			name: "add image",
			images: []string{
				"bryanl/slim-hello-world:v1",
			},
			wanted: []string{"docker.io/bryanl/slim-hello-world:v1"},
		},
		{
			name: "add multiple images",
			images: []string{
				"bryanl/slim-hello-world:v2",
				"bryanl/slim-hello-world:v1",
			},
			wanted: []string{
				"docker.io/bryanl/slim-hello-world:v1",
				"docker.io/bryanl/slim-hello-world:v2",
			},
		},
		{
			name: "with an existing image",
			initial: []string{
				"docker.io/bryanl/slim-hello-world:v2",
			},
			images: []string{
				"bryanl/slim-hello-world:v1",
			},
			wanted: []string{
				"docker.io/bryanl/slim-hello-world:v1",
				"docker.io/bryanl/slim-hello-world:v2",
			},
		},
		{
			name: "duplicate existing image",
			initial: []string{
				"docker.io/bryanl/slim-hello-world:v1",
			},
			images: []string{
				"bryanl/slim-hello-world:v1",
			},
			wanted: []string{
				"docker.io/bryanl/slim-hello-world:v1",
			},
		},
		{
			name:   "no image",
			wanted: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withWorkingDirectory(t, func(wd string) {
				bundleDir := initBundle(t, testHarness, "integration", wd)

				var config sheaf.BundleConfig
				configFile := filepath.Join(bundleDir, "bundle.json")

				if len(tc.initial) > 0 {
					readJSONFile(t, configFile, &config)
					list, err := images.New(tc.initial)
					require.NoError(t, err)
					config.Images = &list
					writeJSONFile(t, configFile, config)
				}

				args := []string{"config", "add-image"}
				for i := range tc.images {
					args = append(args, "-i", tc.images[i])
				}

				err := testHarness.runSheaf(bundleDir, defaultSheafRunSettings, args...)
				require.NoError(t, err)

				readJSONFile(t, configFile, &config)

				var actual []string
				if config.Images != nil {
					actual = config.Images.Strings()
				}

				require.Equal(t, tc.wanted, actual)
			})
		})
	}
}
