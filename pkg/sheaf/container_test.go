/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"path/filepath"
	"testing"

	"github.com/bryanl/sheaf/pkg/images"
	"github.com/stretchr/testify/require"
)

func TestContainers(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantErr  bool
		expected []string
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := filepath.Join("testdata", tt.path)

			got, err := ContainerImages(manifestPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			expected, err := images.New(tt.expected)
			require.NoError(t, err)
			require.Equal(t, expected, got)
		})
	}
}
