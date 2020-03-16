// +build !integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/manifest"
	"github.com/bryanl/sheaf/pkg/sheaf"
	"github.com/pivotal/image-relocation/pkg/images"
)

func TestContainers(t *testing.T) {
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
					JSONPath:   "{.spec.image}",
					Type:       "single",
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
					JSONPath:   "{range .spec.images[*]}{@}{','}{end}",
					Type:       "multiple",
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
					JSONPath:   "{.spec.image}",
					Type:       "invalid",
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
