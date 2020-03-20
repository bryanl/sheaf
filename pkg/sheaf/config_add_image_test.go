/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pivotal/image-relocation/pkg/images"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestConfigAddImage(t *testing.T) {
	genBundleFactory := func(t *testing.T, controller *gomock.Controller, config *mocks.MockBundleConfig) sheaf.BundleFactoryFunc {
		bundle := testutil.GenerateBundle(t, controller,
			testutil.BundleGeneratorConfig(config))
		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	cases := []struct {
		name          string
		images        []string
		bundleFactory bundleFactoryFunc
		configWriter  func(controller *gomock.Controller) *mocks.MockBundleConfigWriter
		wantErr       bool
	}{
		{
			name:   "with images",
			images: []string{"foo", "bar"},
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetImages().Return(&images.Empty)

				expected, err := images.New(
					"docker.io/library/bar",
					"docker.io/library/foo",
				)
				require.NoError(t, err)
				config.EXPECT().SetImages(&expected)

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name:         "no bundle factory",
			configWriter: noopConfigWriter,
			wantErr:      true,
		},
		{
			name: "load bundle error",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return func(string) (bundle sheaf.Bundle, err error) {
					return nil, fmt.Errorf("error")
				}
			},
			configWriter: noopConfigWriter,
			wantErr:      true,
		},
		{
			name: "no bundle config writer",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				return genBundleFactory(t, controller, config)
			},
			wantErr: true,
		},
		{
			name: "write config error",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetImages().Return(&images.Empty)

				config.EXPECT().SetImages(&images.Empty)

				return genBundleFactory(t, controller, config)
			},
			configWriter: errorConfigWriter,
			wantErr:      true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			options := []sheaf.Option{
				sheaf.WithImages(tc.images),
			}

			if tc.bundleFactory != nil {
				options = append(options, sheaf.WithBundleFactory(
					tc.bundleFactory(controller)))
			}

			if tc.configWriter != nil {
				options = append(options, sheaf.WithBundleConfigWriter(
					tc.configWriter(controller)))
			}

			err := sheaf.ConfigAddImage(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func successfulConfigWriter(controller *gomock.Controller) *mocks.MockBundleConfigWriter {
	bcw := mocks.NewMockBundleConfigWriter(controller)
	bcw.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	return bcw
}

func errorConfigWriter(controller *gomock.Controller) *mocks.MockBundleConfigWriter {
	bcw := mocks.NewMockBundleConfigWriter(controller)
	bcw.EXPECT().Write(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
	return bcw
}

func noopConfigWriter(controller *gomock.Controller) *mocks.MockBundleConfigWriter {
	return mocks.NewMockBundleConfigWriter(controller)
}
