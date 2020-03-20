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

func TestManifestAdd(t *testing.T) {
	genBundleFactory := func(t *testing.T, controller *gomock.Controller, ms *mocks.MockManifestService) sheaf.BundleFactoryFunc {
		bundle := testutil.GenerateBundle(t, controller,
			testutil.BundleGeneratorCreateBundle(func(t *testing.T, controller *gomock.Controller, config sheaf.BundleConfig, manifests []sheaf.BundleManifest) *mocks.MockBundle {
				bundle := mocks.NewMockBundle(controller)
				bundle.EXPECT().Config().Return(config).AnyTimes()

				if ms == nil {
					bundle.EXPECT().Manifests().Return(nil, fmt.Errorf("error"))
				} else {
					bundle.EXPECT().Manifests().Return(ms, nil).AnyTimes()
				}
				return bundle
			}))
		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	cases := []struct {
		name          string
		filePaths     []string
		overwrite     bool
		bundleFactory bundleFactoryFunc
		wantErr       bool
	}{
		{
			name:      "with file paths",
			filePaths: []string{"foo", "bar"},
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				ms := mocks.NewMockManifestService(controller)
				ms.EXPECT().Add(false, "foo", "bar").Return(nil)

				return genBundleFactory(t, controller, ms)
			},
		},
		{
			name:      "with file paths with force overwrite",
			filePaths: []string{"foo", "bar"},
			overwrite: true,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				ms := mocks.NewMockManifestService(controller)
				ms.EXPECT().Add(true, "foo", "bar").Return(nil)

				return genBundleFactory(t, controller, ms)
			},
		},
		{
			name:    "with no bundle factory",
			wantErr: true,
		},
		{
			name: "unable to load bundle",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return func(string) (bundle sheaf.Bundle, err error) {
					return nil, fmt.Errorf("error")
				}
			},
			wantErr: true,
		},
		{
			name: "unable to load bundle manifest service",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller, nil)
			},
			wantErr: true,
		},
		{
			name:      "unable to add files",
			filePaths: []string{"foo", "bar"},
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				ms := mocks.NewMockManifestService(controller)
				ms.EXPECT().Add(false, "foo", "bar").Return(fmt.Errorf("error"))

				return genBundleFactory(t, controller, ms)
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			options := []sheaf.Option{
				sheaf.WithFilePaths(tc.filePaths),
				sheaf.WithForce(tc.overwrite),
			}

			if tc.bundleFactory != nil {
				options = append(options, sheaf.WithBundleFactory(
					tc.bundleFactory(controller)))
			}

			err := sheaf.ManifestAdd(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
