/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf_test

import (
	"bytes"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestConfigPush(t *testing.T) {
	genBundleFactory := func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
		bundle := testutil.GenerateBundle(t, controller)
		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	cases := []struct {
		name          string
		ref           string
		bundleFactory bundleFactoryFunc
		bundleImager  func(controller *gomock.Controller) *mocks.MockBundleImager
		imageWriter   func(controller *gomock.Controller) *mocks.MockImageWriter
		wantErr       bool
	}{
		{
			name: "in general",
			ref:  "ref",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(controller)
			},
			bundleImager: func(controller *gomock.Controller) *mocks.MockBundleImager {
				image, err := random.Image(64, 1)
				require.NoError(t, err)

				bi := mocks.NewMockBundleImager(controller)
				bi.EXPECT().
					CreateImage(gomock.Any()).
					Return(image, nil)

				return bi
			},
			imageWriter: func(controller *gomock.Controller) *mocks.MockImageWriter {
				iw := mocks.NewMockImageWriter(controller)
				iw.EXPECT().
					Write("ref", gomock.Any()).Return(nil)

				return iw
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			var buf bytes.Buffer
			r := reporter.New(reporter.WithWriter(&buf))

			options := []sheaf.Option{
				sheaf.WithReference(tc.ref),
				sheaf.WithBundleFactory(tc.bundleFactory(controller)),
				sheaf.WithBundleImager(tc.bundleImager(controller)),
				sheaf.WithImageWriter(tc.imageWriter(controller)),
				sheaf.WithReporter(r),
			}

			err := sheaf.ConfigPush(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
