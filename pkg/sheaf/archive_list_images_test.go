/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestArchiveListImages(t *testing.T) {
	bundleImages := []sheaf.BundleImage{{}}

	genBundleFactory := func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
		is := mocks.NewMockImageService(controller)
		is.EXPECT().List().Return(bundleImages, nil)

		as := mocks.NewMockArtifactsService(controller)
		as.EXPECT().Image().Return(is)

		bundle := testutil.GenerateBundle(t, controller)
		bundle.EXPECT().Artifacts().Return(as)
		bundle.EXPECT().Path().Return("")

		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	tests := []struct {
		name          string
		archive       string
		archiver      func(controller *gomock.Controller) *mocks.MockArchiver
		bundleFactory bundleFactoryFunc
		codec         func(controller *gomock.Controller) *mocks.MockCodec
		wantErr       bool
	}{
		{
			name:    "in general",
			archive: "archive.tgz",
			archiver: func(controller *gomock.Controller) *mocks.MockArchiver {
				a := mocks.NewMockArchiver(controller)
				a.EXPECT().
					UnarchivePath("archive.tgz", gomock.Any()).
					Return(nil)

				return a
			},
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(controller)
			},
			codec: func(controller *gomock.Controller) *mocks.MockCodec {
				codec := mocks.NewMockCodec(controller)
				codec.EXPECT().Encode(gomock.Any()).Return([]byte("data"), nil)
				return codec
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			options := []sheaf.Option{
				sheaf.WithArchive(test.archive),
			}

			if test.archiver != nil {
				options = append(options, sheaf.WithArchiver(test.archiver(controller)))
			}

			if test.bundleFactory != nil {
				options = append(options, sheaf.WithBundleFactory(test.bundleFactory(controller)))
			}

			if test.codec != nil {
				options = append(options, sheaf.WithCodec(test.codec(controller)))
			}

			err := sheaf.ArchiveListImages(options...)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
