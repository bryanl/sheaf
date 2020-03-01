/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

type listMocks struct {
	artifacts *mocks.MockArtifactsService
	decoder   *mocks.MockDecoder
}

func TestImageService_List(t *testing.T) {
	cases := []struct {
		name    string
		init    func(t *testing.T, m listMocks)
		wanted  []sheaf.BundleImage
		wantErr bool
	}{
		{
			name: "in general",
			init: func(t *testing.T, m listMocks) {
				data, err := ioutil.ReadFile(filepath.Join("testdata", "index.json"))
				require.NoError(t, err)
				m.artifacts.EXPECT().Index().Return(data, nil)

				m.decoder.EXPECT().Decode(gomock.Any(), gomock.Any()).
					DoAndReturn(json.Unmarshal)
			},
			wanted: []sheaf.BundleImage{
				{
					Name:      "example/image-a:abc",
					Digest:    "sha256:4528b0a54dd4ec91f0398856216b24532566618340c7ef6fd00345b776fb2c10",
					MediaType: "application/vnd.docker.distribution.manifest.v2+json",
				},
				{
					Name:      "example/image-b:v1.0.0",
					Digest:    "sha256:a7358e4600ae00bf976dba9c299c3dcd7bd0473e18ff334dde35ba0f6535663b",
					MediaType: "application/vnd.docker.distribution.manifest.v2+json",
				},
			},
		},
		{
			name: "image is missing ref name",
			init: func(t *testing.T, m listMocks) {
				data, err := ioutil.ReadFile(filepath.Join("testdata", "index-no-ref-name.json"))
				require.NoError(t, err)
				m.artifacts.EXPECT().Index().Return(data, nil)

				m.decoder.EXPECT().Decode(gomock.Any(), gomock.Any()).
					DoAndReturn(json.Unmarshal)
			},
		},
		{
			name: "index file error",
			init: func(t *testing.T, m listMocks) {
				m.artifacts.EXPECT().Index().Return(nil, fmt.Errorf("error"))
			},
			wantErr: true,
		},
		{
			name: "index data is invalid",
			init: func(t *testing.T, m listMocks) {
				data, err := ioutil.ReadFile(filepath.Join("testdata", "index.json"))
				require.NoError(t, err)
				m.artifacts.EXPECT().Index().Return(data, nil)

				m.decoder.EXPECT().Decode(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("decode error"))
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			artifactsService := mocks.NewMockArtifactsService(controller)
			decoder := mocks.NewMockDecoder(controller)

			if tc.init != nil {
				tc.init(t, listMocks{
					artifacts: artifactsService,
					decoder:   decoder,
				})
			}

			imageService := NewImageService(artifactsService, ImageServiceDecoder(decoder))

			actual, err := imageService.List()
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.wanted, actual)

		})
	}

}
