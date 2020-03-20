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

	"github.com/pivotal/image-relocation/pkg/images"

	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// StageFile stages a file to a destination.
func StageFile(t *testing.T, name, dest string) []byte {
	data, err := ioutil.ReadFile(filepath.Join("testdata", name))
	require.NoError(t, err)
	require.NoError(t, ioutil.WriteFile(dest, data, 0600))

	return data
}

// Testdata returns test data as bytes.
func Testdata(t *testing.T, parts ...string) []byte {
	data, err := ioutil.ReadFile(filepath.Join(append([]string{"testdata"}, parts...)...))
	require.NoError(t, err)

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

// BundleGeneratorManifests sets the bundle manifests for BundleGenerator.
func BundleGeneratorManifests(list []sheaf.BundleManifest) BundleGeneratorOption {
	return func(generator BundleGenerator) BundleGenerator {
		generator.manifests = list
		return generator
	}
}

// CreateBundleFunc is a function that create a Bundle mock.
type CreateBundleFunc func(t *testing.T, controller *gomock.Controller, config sheaf.BundleConfig, manifests []sheaf.BundleManifest) *mocks.MockBundle

// BundleGeneratorCreateBundle sets the bundle creator for BundleGenerator.
func BundleGeneratorCreateBundle(fn CreateBundleFunc) BundleGeneratorOption {
	return func(generator BundleGenerator) BundleGenerator {
		generator.createBundleFunc = fn
		return generator
	}
}

// BundleGenerator generates a sheaf.Bundle mock.
type BundleGenerator struct {
	config           sheaf.BundleConfig
	manifests        []sheaf.BundleManifest
	createBundleFunc CreateBundleFunc
}

// DefaultBundleStubs are the default bundle stubs.
func DefaultBundleStubs(t *testing.T, controller *gomock.Controller, config sheaf.BundleConfig, manifests []sheaf.BundleManifest) *mocks.MockBundle {
	bundle := mocks.NewMockBundle(controller)
	bundle.EXPECT().Config().Return(config).AnyTimes()

	m := mocks.NewMockManifestService(controller)
	m.EXPECT().List().Return(manifests, nil).AnyTimes()

	bundle.EXPECT().Manifests().Return(m, nil).AnyTimes()

	imageList, err := images.New("image")
	require.NoError(t, err)
	bundle.EXPECT().Images().Return(imageList, nil).AnyTimes()

	return bundle
}

// GenerateBundleConfig generates a bundle config mock.
func GenerateBundleConfig(controller *gomock.Controller) *mocks.MockBundleConfig {
	bc := mocks.NewMockBundleConfig(controller)
	bc.EXPECT().GetName().Return("project").AnyTimes()
	bc.EXPECT().GetVersion().Return("0.1.0").AnyTimes()
	bc.EXPECT().GetSchemaVersion().Return("v1alpha1").AnyTimes()

	return bc
}

// GenerateBundle generates a bundle mock.
func GenerateBundle(t *testing.T, controller *gomock.Controller, options ...BundleGeneratorOption) *mocks.MockBundle {
	bg := BundleGenerator{
		config: GenerateBundleConfig(controller),
		manifests: []sheaf.BundleManifest{
			{
				ID:   "deploy.yaml",
				Data: sampleManifests,
			},
		},
	}

	for _, option := range options {
		bg = option(bg)
	}

	fn := bg.createBundleFunc
	if fn == nil {
		fn = DefaultBundleStubs
	}

	return fn(t, controller, bg.config, bg.manifests)
}

var sampleManifests = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx-1
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-1
  template:
    metadata:
      labels:
        app: nginx-1
    spec:
      containers:
        - name: nginx
          image: nginx:1.17.8
          ports:
            - containerPort: 80
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-2
  labels:
    app: nginx-2
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx-2
  template:
    metadata:
      labels:
        app: nginx-2
    spec:
      containers:
        - name: nginx
          image: nginx:1.17.8
          ports:
            - containerPort: 80`)
