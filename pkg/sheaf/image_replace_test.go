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
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceImage(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		oldImage     string
		newImage     string
		expectedPath string
	}{
		{
			name:         "deployment",
			path:         "deployment.yaml",
			oldImage:     "nginx:1.7.9",
			newImage:     "example.com/nginx:1.7.9",
			expectedPath: "deployment-replaced.yaml",
		},
		{
			name:         "quoted",
			path:         "quoted.yaml",
			oldImage:     "quay.io/jetstack/cert-manager-cainjector@sha256:9ff6923f6c567573103816796df283d03256bc7a9edb7450542e106b349cf34a",
			newImage:     "example.com/jetstack/cert-manager-cainjector@sha256:9ff6923f6c567573103816796df283d03256bc7a9edb7450542e106b349cf34a",
			expectedPath: "quoted-replaced.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := readTestData(tt.path, t)
			updatedManifest := string(replaceImage(manifest, tt.oldImage, tt.newImage))
			expectedManifest := string(readTestData(tt.expectedPath, t))

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
