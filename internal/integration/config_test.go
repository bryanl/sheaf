// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"bytes"
	"io/ioutil"
	"os"
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
			withWorkingDirectory(t, func(options wdOptions) {
				b := sheafInit(t, testHarness, "integration", options.dir)

				b.updateConfig(func(config *sheaf.BundleConfig) {
					list, err := images.New(tc.initial)
					require.NoError(t, err)
					config.Images = &list
				})

				args := []string{"config", "add-image"}
				for i := range tc.images {
					args = append(args, "-i", tc.images[i])
				}

				err := b.harness.runSheaf(b.dir, defaultSheafRunSettings, args...)
				require.NoError(t, err)

				config := b.readConfig()

				var actual []string
				if config.Images != nil {
					actual = config.Images.Strings()
				}

				require.Equal(t, tc.wanted, actual)
			})
		})
	}
}

func Test_sheaf_config_set_udi(t *testing.T) {

	cases := []struct {
		name     string
		existing []sheaf.UserDefinedImage
		udi      udi
		wanted   []sheaf.UserDefinedImage
	}{
		{
			name: "set user defined image",
			udi: udi{
				APIVersion: "example.com/v1",
				Kind:       "Resource",
				JSONPath:   "{.spec.image}",
			},
			wanted: []sheaf.UserDefinedImage{
				{
					APIVersion: "example.com/v1",
					Kind:       "Resource",
					JSONPath:   "{.spec.image}",
					Type:       sheaf.SingleResult,
				},
			},
		},
		{
			name: "set user defined image with existing",
			existing: []sheaf.UserDefinedImage{
				{
					APIVersion: "example.com/v2",
					Kind:       "Resource",
					JSONPath:   "{.spec.image}",
					Type:       sheaf.SingleResult,
				},
			},
			udi: udi{
				APIVersion: "example.com/v1",
				Kind:       "Resource",
				JSONPath:   "{.spec.image}",
			},
			wanted: []sheaf.UserDefinedImage{
				{
					APIVersion: "example.com/v1",
					Kind:       "Resource",
					JSONPath:   "{.spec.image}",
					Type:       sheaf.SingleResult,
				},
				{
					APIVersion: "example.com/v2",
					Kind:       "Resource",
					JSONPath:   "{.spec.image}",
					Type:       sheaf.SingleResult,
				},
			},
		},
		{
			name: "update existing image",
			existing: []sheaf.UserDefinedImage{
				{
					APIVersion: "example.com/v1",
					Kind:       "Resource",
					JSONPath:   "{.spec.image}",
					Type:       sheaf.SingleResult,
				},
			},
			udi: udi{
				APIVersion: "example.com/v1",
				Kind:       "Resource",
				JSONPath:   "{.spec.images}",
				Type:       "multiple",
			},
			wanted: []sheaf.UserDefinedImage{
				{
					APIVersion: "example.com/v1",
					Kind:       "Resource",
					JSONPath:   "{.spec.images}",
					Type:       sheaf.MultiResult,
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withWorkingDirectory(t, func(options wdOptions) {
				b := sheafInit(t, testHarness, "integration", options.dir)

				b.updateConfig(func(config *sheaf.BundleConfig) {
					config.UserDefinedImages = tc.existing
				})

				args := append([]string{"config", "set-udi"}, tc.udi.toArgs()...)

				err := testHarness.runSheaf(b.dir, defaultSheafRunSettings, args...)
				require.NoError(t, err)

				config := b.readConfig()

				require.Equal(t, tc.wanted, config.UserDefinedImages)
			})
		})
	}
}

func Test_sheaf_config_delete_udi(t *testing.T) {
	cases := []struct {
		name     string
		existing []sheaf.UserDefinedImage
		id       udiID
		wanted   []sheaf.UserDefinedImage
	}{
		{
			name: "delete user defined image",
			existing: []sheaf.UserDefinedImage{
				{
					APIVersion: "example.com/v1",
					Kind:       "Resource",
					JSONPath:   "{.spec.image}",
					Type:       sheaf.SingleResult,
				},
			},
			id: udiID{
				APIVersion: "example.com/v1",
				Kind:       "Resource",
			},
			wanted: nil,
		},
		{
			name: "delete user defined image that does not exist",
			id: udiID{
				APIVersion: "example.com/v1",
				Kind:       "Resource",
			},
			wanted: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withWorkingDirectory(t, func(options wdOptions) {
				b := sheafInit(t, testHarness, "integration", options.dir)
				b.updateConfig(func(config *sheaf.BundleConfig) {
					config.UserDefinedImages = tc.existing
				})

				args := append([]string{"config", "delete-udi"}, tc.id.toArgs()...)

				err := testHarness.runSheaf(b.dir, defaultSheafRunSettings, args...)
				require.NoError(t, err)

				config := b.readConfig()
				require.Equal(t, tc.wanted, config.UserDefinedImages)
			})
		})
	}
}

func Test_sheaf_config_get(t *testing.T) {
	withWorkingDirectory(t, func(options wdOptions) {
		b := sheafInit(t, testHarness, "integration", options.dir)

		settings := genSheafRunSettings()
		var actual bytes.Buffer
		settings.Stdout = &actual

		args := append([]string{"config", "get"})

		err := b.harness.runSheaf(b.dir, settings, args...)
		require.NoError(t, err)

		wanted := readFile(t, b.configFile())
		wanted = bytes.TrimSpace(wanted)

		require.Equal(t, string(wanted), actual.String())
	})
}

func Test_sheaf_config_push_and_pull(t *testing.T) {
	withWorkingDirectory(t, func(options wdOptions) {

		b := sheafInit(t, testHarness, "integration", options.dir)

		settings := defaultSheafRunSettings

		require.NoError(t, b.harness.runSheaf(b.dir, settings, "manifest", "add",
			"-f", testdata(t, "config", "push")))

		ref := genRegistryPath(options)

		pushArgs := append([]string{"config", "push", b.dir, ref, "--insecure-registry"})
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

type udiID struct {
	APIVersion string
	Kind       string
}

func (u udiID) toArgs() []string {
	args := []string{
		"--api-version", u.APIVersion,
		"--kind", u.Kind,
	}

	return args
}

type udi struct {
	APIVersion string
	Kind       string
	JSONPath   string
	Type       string
}

func (u udi) toArgs() []string {
	args := []string{
		"--api-version", u.APIVersion,
		"--kind", u.Kind,
		"--json-path", u.JSONPath,
	}

	if u.Type != "" {
		args = append(args, "--type", u.Type)
	}

	return args
}
