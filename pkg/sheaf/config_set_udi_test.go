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

func TestConfigSetUDI(t *testing.T) {
	genBundleFactory := func(t *testing.T, controller *gomock.Controller, config *mocks.MockBundleConfig) sheaf.BundleFactoryFunc {
		bundle := testutil.GenerateBundle(t, controller,
			testutil.BundleGeneratorConfig(config))
		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	udi1 := sheaf.UserDefinedImage{
		APIVersion: "v1",
		Kind:       "Kind1",
		JSONPath:   "{.}",
		Type:       sheaf.SingleResult,
	}
	udi2 := sheaf.UserDefinedImage{
		APIVersion: "v2",
		Kind:       "Kind1",
		JSONPath:   "{.}",
		Type:       sheaf.SingleResult,
	}
	udi3 := sheaf.UserDefinedImage{
		APIVersion: "v1",
		Kind:       "Kind2",
		JSONPath:   "{.}",
		Type:       sheaf.SingleResult,
	}

	cases := []struct {
		name          string
		udi           sheaf.UserDefinedImage
		bundleFactory bundleFactoryFunc
		configWriter  func(controller *gomock.Controller) *mocks.MockBundleConfigWriter
		wantErr       bool
	}{
		{
			name: "with no existing UDIs",
			udi:  udi1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return(nil)
				config.EXPECT().SetUserDefinedImages([]sheaf.UserDefinedImage{udi1})

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name: "with existing UDIs",
			udi:  udi1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{udi2})
				config.EXPECT().SetUserDefinedImages([]sheaf.UserDefinedImage{udi1, udi2})

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name: "writes UDIs in a stable order (sorts api version)",
			udi:  udi2,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{udi1})
				config.EXPECT().SetUserDefinedImages([]sheaf.UserDefinedImage{udi1, udi2})

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name: "writes UDIs in a stable order (sorts kind)",
			udi:  udi1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{udi3})
				config.EXPECT().SetUserDefinedImages([]sheaf.UserDefinedImage{udi1, udi3})

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name: "update existing udi",
			udi:  udi1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{udi1})
				config.EXPECT().SetUserDefinedImages([]sheaf.UserDefinedImage{udi1})

				return genBundleFactory(t, controller, config)
			},
			configWriter: successfulConfigWriter,
		},
		{
			name:         "with no bundle factory",
			udi:          udi1,
			configWriter: noopConfigWriter,
			wantErr:      true,
		},
		{
			name: "no config writer",
			udi:  udi1,
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
			name: "write config error",
			udi:  udi1,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				config := testutil.GenerateBundleConfig(controller)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{udi1})
				config.EXPECT().SetUserDefinedImages([]sheaf.UserDefinedImage{udi1})

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
				sheaf.WithUserDefinedImage(tc.udi),
			}

			if tc.bundleFactory != nil {
				options = append(options, sheaf.WithBundleFactory(
					tc.bundleFactory(controller)))
			}

			if tc.configWriter != nil {
				options = append(options, sheaf.WithBundleConfigWriter(
					tc.configWriter(controller)))
			}

			err := sheaf.ConfigSetUDI(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
