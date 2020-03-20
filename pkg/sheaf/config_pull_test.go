/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf_test

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/fake"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestConfigPull(t *testing.T) {
	cases := []struct {
		name    string
		options func(controller *gomock.Controller) []sheaf.Option
		wantErr bool
	}{
		{
			name: "in general",
			options: func(controller *gomock.Controller) []sheaf.Option {
				archiver := mocks.NewMockArchiver(controller)
				archiver.EXPECT().Unarchive(gomock.Any(), "dest").Return(nil)

				imageReader := mocks.NewMockImageReader(controller)
				image, err := random.Image(100, 1)
				require.NoError(t, err)
				imageReader.EXPECT().
					Read("ref").Return(image, nil)

				return []sheaf.Option{
					sheaf.WithReference("ref"),
					sheaf.WithArchiver(archiver),
					sheaf.WithDestination("dest"),
					sheaf.WithImageReader(imageReader),
				}
			},
		},
		{
			name: "image reader failed",
			options: func(controller *gomock.Controller) []sheaf.Option {
				imageReader := mocks.NewMockImageReader(controller)
				imageReader.EXPECT().Read("ref").Return(nil, fmt.Errorf("error"))
				return []sheaf.Option{
					sheaf.WithReference("ref"),
					sheaf.WithDestination("dest"),
					sheaf.WithImageReader(imageReader),
				}
			},
			wantErr: true,
		},
		{
			name: "image layers failed",
			options: func(controller *gomock.Controller) []sheaf.Option {
				imageReader := mocks.NewMockImageReader(controller)

				image := &fake.FakeImage{}
				image.LayersReturns(nil, fmt.Errorf("error"))
				imageReader.EXPECT().
					Read("ref").Return(image, nil)

				return []sheaf.Option{
					sheaf.WithReference("ref"),
					sheaf.WithDestination("dest"),
					sheaf.WithImageReader(imageReader),
				}
			},
			wantErr: true,
		},
		{
			name: "invalid image format (wrong layer count)",
			options: func(controller *gomock.Controller) []sheaf.Option {
				imageReader := mocks.NewMockImageReader(controller)
				image, err := random.Image(100, 2)
				require.NoError(t, err)
				imageReader.EXPECT().
					Read("ref").Return(image, nil)

				return []sheaf.Option{
					sheaf.WithReference("ref"),
					sheaf.WithDestination("dest"),
					sheaf.WithImageReader(imageReader),
				}
			},
			wantErr: true,
		},
		{
			name: "unable to retrieve compressed layer",
			options: func(controller *gomock.Controller) []sheaf.Option {
				imageReader := mocks.NewMockImageReader(controller)

				layer := &fakeLayer{}

				image := &fake.FakeImage{}
				image.LayersReturns([]v1.Layer{layer}, nil)
				imageReader.EXPECT().
					Read("ref").Return(image, nil)

				return []sheaf.Option{
					sheaf.WithReference("ref"),
					sheaf.WithDestination("dest"),
					sheaf.WithImageReader(imageReader),
				}
			},
			wantErr: true,
		},
		{
			name: "unarchive fails",
			options: func(controller *gomock.Controller) []sheaf.Option {
				archiver := mocks.NewMockArchiver(controller)
				archiver.EXPECT().Unarchive(gomock.Any(), "dest").Return(fmt.Errorf("error"))

				imageReader := mocks.NewMockImageReader(controller)
				image, err := random.Image(100, 1)
				require.NoError(t, err)
				imageReader.EXPECT().
					Read("ref").Return(image, nil)

				return []sheaf.Option{
					sheaf.WithReference("ref"),
					sheaf.WithArchiver(archiver),
					sheaf.WithDestination("dest"),
					sheaf.WithImageReader(imageReader),
				}
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			var buf bytes.Buffer
			options := []sheaf.Option{
				sheaf.WithReporter(reporter.New(reporter.WithWriter(&buf))),
			}
			if tc.options != nil {
				options = append(options, tc.options(controller)...)
			}

			err := sheaf.ConfigPull(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

type fakeLayer struct{ v1.Layer }

func (l *fakeLayer) Compressed() (io.ReadCloser, error) {
	return nil, fmt.Errorf("error")
}
