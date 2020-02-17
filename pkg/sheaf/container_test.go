/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"path/filepath"
	"testing"

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
			require.Equal(t, tt.expected, got)
		})
	}
}
