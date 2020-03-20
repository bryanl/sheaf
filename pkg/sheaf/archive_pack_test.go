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

func TestArchivePack(t *testing.T) {
	genBundleFactory := func(t *testing.T, controller *gomock.Controller) sheaf.BundleFactoryFunc {
		bundle := testutil.GenerateBundle(t, controller)
		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	cases := []struct {
		name          string
		destination   string
		force         bool
		bundleFactory bundleFactoryFunc
		bundlePacker  func(controller *gomock.Controller) *mocks.MockBundlePacker
		wantErr       bool
	}{
		{
			name:        "in general",
			destination: "dest",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller)
			},
			bundlePacker: func(controller *gomock.Controller) *mocks.MockBundlePacker {
				bp := mocks.NewMockBundlePacker(controller)
				bp.EXPECT().
					Pack(gomock.Any(), "dest", false).
					Return(nil)
				return bp
			},
		},
		{
			name:        "force overwrite",
			destination: "dest",
			force:       true,
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller)
			},
			bundlePacker: func(controller *gomock.Controller) *mocks.MockBundlePacker {
				bp := mocks.NewMockBundlePacker(controller)
				bp.EXPECT().
					Pack(gomock.Any(), "dest", true).
					Return(nil)
				return bp
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			options := []sheaf.Option{
				sheaf.WithBundleFactory(tc.bundleFactory(controller)),
				sheaf.WithBundlePacker(tc.bundlePacker(controller)),
				sheaf.WithDestination(tc.destination),
				sheaf.WithForce(tc.force),
			}

			err := sheaf.ArchivePack(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
