/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package testutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
	"github.com/pivotal/image-relocation/pkg/images"
)

// StageFile stages a file to a destination.
func StageFile(t *testing.T, name, dest string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	require.NoError(t, ioutil.WriteFile(dest, data, 0600))

	return data
}

// WithBundleDir creates a bundle directory and runs a function.
func WithBundleDir(t *testing.T, fn func(dir string)) {
	dir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(dir))
	}()

	fn(dir)
}

// SlurpData reads bytes from a file.
func SlurpData(t *testing.T, source string) []byte {
	data, err := ioutil.ReadFile(source)
	require.NoError(t, err)
	return data
}

// BundleGeneratorOption is a functional option for configuring BundleGenerator.
type BundleGeneratorOption func(generator BundleGenerator) BundleGenerator

// BundleGeneratorConfig sets config for BundleGenerator.
func BundleGeneratorConfig(config sheaf.BundleConfig) BundleGeneratorOption {
	return func(generator BundleGenerator) BundleGenerator {
		generator.config = config
		return generator
	}
}

// BundleGenerator generates a sheaf.Bundle mock.
type BundleGenerator struct {
	config    sheaf.BundleConfig
	manifests []sheaf.BundleManifest
}

var (
	// BundleConfig is the default bundle config for mocks.
	BundleConfig = sheaf.BundleConfig{
		Name:          "project",
		Version:       "0.1.0",
		SchemaVersion: "v1alpha1",
	}
)

// GenerateBundle generates a bundle mock.
func GenerateBundle(t *testing.T, controller *gomock.Controller, options ...BundleGeneratorOption) *mocks.MockBundle {
	bg := BundleGenerator{
		config: BundleConfig,
		manifests: []sheaf.BundleManifest{
			{
				ID:   "deploy.yaml",
				Data: SlurpData(t, filepath.Join("testdata", "manifests", "deploy.yaml")),
			},
		},
	}

	for _, option := range options {
		bg = option(bg)
	}

	bundle := mocks.NewMockBundle(controller)
	bundle.EXPECT().Config().Return(bg.config).AnyTimes()

	m := mocks.NewMockManifestService(controller)
	m.EXPECT().List().Return(bg.manifests, nil).AnyTimes()

	bundle.EXPECT().Manifests().Return(m, nil).AnyTimes()

	imageList, err := images.New("image")
	require.NoError(t, err)
	bundle.EXPECT().Images().Return(imageList, nil).AnyTimes()

	return bundle
}
