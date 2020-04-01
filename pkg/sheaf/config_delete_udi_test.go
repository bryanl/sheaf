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
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestConfigDeleteUDI(t *testing.T) {
	genBundleFactory := func(t *testing.T, controller *gomock.Controller, config *mocks.MockBundleConfig) sheaf.BundleFactoryFunc {
		bundle := testutil.GenerateBundle(t, controller,
			testutil.BundleGeneratorConfig(config))
		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	udiKey1 := sheaf.UserDefinedImageKey{
		APIVersion: "v1",
		Kind:       "Kind1",
	}

	udi1 := sheaf.UserDefinedImage{
		APIVersion: "v1",
		Kind:       "Kind1",
		JSONPath:   "{.}",
	}

	cases := []struct {
		name          string
		udiKey        sheaf.UserDefinedImageKey
		bundleFactory bundleFactoryFunc
		configWriter  func(controller *gomock.Controller) *mocks.MockBundleConfigWriter
		wantErr       bool
	}{
		{
			name:   "with no existing UDIs",
			udiKey: udiKey1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return(nil)
				config.EXPECT().SetUserDefinedImages(nil)

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name:   "with existing UDIs",
			udiKey: udiKey1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{udi1})
				config.EXPECT().SetUserDefinedImages(nil)

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name:         "with no bundle factory",
			udiKey:       udiKey1,
			configWriter: noopConfigWriter,
			wantErr:      true,
		},
		{
			name:   "no config writer",
			udiKey: udiKey1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				return genBundleFactory(t, controller, config)
			},
			wantErr: true,
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
			name:   "write config error",
			udiKey: udiKey1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{udi1})
				config.EXPECT().SetUserDefinedImages(nil)

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
				sheaf.WithUserDefinedImageKey(tc.udiKey),
			}

			if tc.bundleFactory != nil {
				options = append(options, sheaf.WithBundleFactory(
					tc.bundleFactory(controller)))
			}

			if tc.configWriter != nil {
				options = append(options, sheaf.WithBundleConfigWriter(
					tc.configWriter(controller)))
			}

			err := sheaf.ConfigDeleteUDI(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
