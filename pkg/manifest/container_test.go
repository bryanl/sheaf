// +build !integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest_test

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/manifest"
	"github.com/bryanl/sheaf/pkg/sheaf"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/images"
	"github.com/stretchr/testify/require"
)

func TestContainerImages(t *testing.T) {
	tests := []struct {
		name              string
		path              string
		userDefinedImages []sheaf.UserDefinedImage
		wantErr           bool
		expected          []string
	}{
		{
			name:    "invalid file",
			path:    "invalid",
			wantErr: true,
		},
		{
			name: "deployment",
			path: "deployment.yaml",
			expected: []string{
				"nginx:1.7.9",
			},
		},
		{
			name: "pod",
			path: "pod.yaml",
			expected: []string{
				"busybox",
			},
		},
		{
			name:     "service",
			path:     "service.yaml",
			expected: nil,
		},
		{
			name: "multi",
			path: "multi.yaml",
			expected: []string{
				"nginx:1.17.8",
			},
		},
		{
			name: "quoted",
			path: "quoted.yaml",
			expected: []string{
				"quay.io/jetstack/cert-manager-cainjector@sha256:9ff6923f6c567573103816796df283d03256bc7a9edb7450542e106b349cf34a",
			},
		},
		{
			name: "user defined: single",
			path: "user-defined-single.yaml",
			userDefinedImages: []sheaf.UserDefinedImage{
				{
					APIVersion: "caching.internal.knative.dev/v1alpha1",
					Kind:       "Image",
					JSONPath:   ".spec.image",
				},
			},
			expected: []string{
				"gcr.io/example/image1",
			},
		},
		{
			name: "user defined: multi",
			path: "user-defined-multi.yaml",
			userDefinedImages: []sheaf.UserDefinedImage{
				{
					APIVersion: "example.dev/v1",
					Kind:       "Foo",
					JSONPath:   ".spec.images[*]",
				},
			},
			expected: []string{
				"gcr.io/example/image1",
				"gcr.io/example/image2",
			},
		},
		{
			name: "user defined: error",
			path: "user-defined-single.yaml",
			userDefinedImages: []sheaf.UserDefinedImage{
				{
					APIVersion: "caching.internal.knative.dev/v1alpha1",
					Kind:       "Image",
					JSONPath:   "]",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := filepath.Join("testdata", tt.path)

			got, err := manifest.ContainerImages(manifestPath, tt.userDefinedImages)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected, err := images.New(tt.expected...)
			require.NoError(t, err)
			require.Equal(t, expected, got)
		})
	}
}

func TestMapContainer(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		mapping      map[string]string
		expectedPath string
	}{
		{
			name:         "deployment",
			path:         "deployment.yaml",
			mapping:      map[string]string{"nginx:1.7.9": "example.com/nginx:1.7.9"},
			expectedPath: "deployment-replaced.yaml",
		},
		{
			name:         "synonym",
			path:         "deployment-synonym.yaml",
			mapping:      map[string]string{"nginx:1.7.9": "example.com/nginx:1.7.9"},
			expectedPath: "deployment-replaced.yaml",
		},
		{
			name:         "quoted",
			path:         "quoted.yaml",
			mapping:      map[string]string{"quay.io/jetstack/cert-manager-cainjector@sha256:9ff6923f6c567573103816796df283d03256bc7a9edb7450542e106b349cf34a": "example.com/jetstack/cert-manager-cainjector@sha256:9ff6923f6c567573103816796df283d03256bc7a9edb7450542e106b349cf34a"},
			expectedPath: "quoted-replaced.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			man := readTestData(tt.path, t)

			mapping := map[image.Name]image.Name{}
			for old, new := range tt.mapping {
				oldName, err := image.NewName(old)
				require.NoError(t, err)

				newName, err := image.NewName(new)
				require.NoError(t, err)

				mapping[oldName] = newName
			}

			newMan, err := manifest.MapContainer(man, nil, func(originalImage image.Name) (image.Name, error) {
				newImage, ok := mapping[originalImage]
				require.True(t, ok)
				return newImage, nil
			})
			require.NoError(t, err)

			updatedManifest := string(testutil.NormalizeNewlines(newMan))
			expectedManifest := string(testutil.NormalizeNewlines(readTestData(tt.expectedPath, t)))

			require.Equal(t, expectedManifest, updatedManifest)
		})
	}
}

func readTestData(filename string, t *testing.T) []byte {
	path := filepath.Join("testdata", filename)
	data, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	return data
}
