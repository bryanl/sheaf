/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/stretchr/testify/require"
)

func TestReplaceImage(t *testing.T) {
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
		{
			name:         "non-standard",
			path:         "non-standard.yaml",
			mapping:      map[string]string{"gcr.io/cf-build-service-public/kpack/build-init@sha256:5205844aefba7c91803198ef81da9134031f637d605d293dfe4531c622aa42b1": "example.com/cf-build-service-public/kpack/build-init@sha256:5205844aefba7c91803198ef81da9134031f637d605d293dfe4531c622aa42b1"},
			expectedPath: "non-standard-replaced.yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifest := readTestData(tt.path, t)

			mapping := map[image.Name]image.Name{}
			for old, new := range tt.mapping {
				oldName, err := image.NewName(old)
				require.NoError(t, err)

				newName, err := image.NewName(new)
				require.NoError(t, err)

				mapping[oldName] = newName
			}

			updatedManifest := string(replaceImage(manifest, mapping))
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
