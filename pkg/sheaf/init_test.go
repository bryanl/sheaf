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
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestInit(t *testing.T) {
	successCreate := func() sheaf.BundleCreatorFunc {
		return func(bc sheaf.BundleConfig) error {
			return nil
		}
	}
	failureCreate := func() sheaf.BundleCreatorFunc {
		return func(bc sheaf.BundleConfig) error {
			return fmt.Errorf("error")
		}
	}

	cases := []struct {
		name                    string
		bundleName              string
		version                 string
		initBundleConfigFactory func(controller *gomock.Controller) sheaf.BundleConfigFactory
		initBundleCreator       func() sheaf.BundleCreatorFunc
		wantErr                 bool
	}{
		{
			name:       "in general",
			bundleName: "project",
			initBundleConfigFactory: func(controller *gomock.Controller) sheaf.BundleConfigFactory {
				bc := testutil.GenerateBundleConfig(controller)
				bc.EXPECT().SetVersion("0.1.0")
				bc.EXPECT().SetName("project")

				return func() sheaf.BundleConfig {
					return bc
				}
			},
			initBundleCreator: successCreate,
		},
		{
			name:       "supply version",
			bundleName: "project",
			version:    "1.0.0",
			initBundleConfigFactory: func(controller *gomock.Controller) sheaf.BundleConfigFactory {
				bc := testutil.GenerateBundleConfig(controller)
				bc.EXPECT().SetVersion("1.0.0")
				bc.EXPECT().SetName("project")

				return func() sheaf.BundleConfig {
					return bc
				}
			},
			initBundleCreator: successCreate,
		},
		{
			name: "no bundle name",
			initBundleConfigFactory: func(controller *gomock.Controller) sheaf.BundleConfigFactory {
				bc := testutil.GenerateBundleConfig(controller)
				return func() sheaf.BundleConfig {
					return bc
				}
			},
			initBundleCreator: successCreate,
			wantErr:           true,
		},
		{
			name:              "no bundle config factory",
			bundleName:        "project",
			initBundleCreator: successCreate,
			wantErr:           true,
		},
		{
			name:       "no bundle creator",
			bundleName: "project",
			initBundleConfigFactory: func(controller *gomock.Controller) sheaf.BundleConfigFactory {
				bc := testutil.GenerateBundleConfig(controller)
				return func() sheaf.BundleConfig {
					return bc
				}
			},
			wantErr: true,
		},
		{
			name:       "create bundle returns error",
			bundleName: "project",
			initBundleConfigFactory: func(controller *gomock.Controller) sheaf.BundleConfigFactory {
				bc := testutil.GenerateBundleConfig(controller)
				bc.EXPECT().SetVersion("0.1.0")
				bc.EXPECT().SetName("project")

				return func() sheaf.BundleConfig {
					return bc
				}
			},
			initBundleCreator: failureCreate,
			wantErr:           true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			options := []sheaf.Option{
				sheaf.WithBundleName(tc.bundleName),
			}
			if tc.initBundleConfigFactory != nil {
				options = append(options, sheaf.WithBundleConfigFactory(tc.initBundleConfigFactory(controller)))
			}

			if tc.version != "" {
				options = append(options, sheaf.WithBundleVersion(tc.version))
			}

			if tc.initBundleCreator != nil {
				options = append(options, sheaf.WithBundleCreator(tc.initBundleCreator()))
			}

			err := sheaf.Init(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}

}
