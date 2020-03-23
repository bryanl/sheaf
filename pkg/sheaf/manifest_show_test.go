/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestManifestShow(t *testing.T) {
	genBundleFactory := func(t *testing.T, controller *gomock.Controller, names []string) sheaf.BundleFactoryFunc {
		var manifests []sheaf.BundleManifest
		for _, name := range names {
			manifests = append(manifests, genManifest(name))
		}
		bundle := testutil.GenerateBundle(t, controller,
			testutil.BundleGeneratorManifests(manifests))
		return func(string) (sheaf.Bundle, error) {
			return bundle, nil
		}
	}

	expectBundleManifest := func(name, wantedPrefix string, ir *mocks.MockImageReplacer) {
		ir.EXPECT().
			Replace(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(manifest sheaf.BundleManifest, config sheaf.BundleConfig, actualPrefix string) ([]byte, error) {
				require.Equal(t, name, manifest.ID)
				require.Equal(t, wantedPrefix, actualPrefix)
				return manifest.Data, nil
			})
	}

	cases := []struct {
		name           string
		bundleFactory  bundleFactoryFunc
		imageRelocator func(controller *gomock.Controller) sheaf.ImageReplacer
		prefix         string
		wantErr        bool
		want           string
	}{
		{
			name: "in general",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller, []string{"deploy1.yaml"})
			},
			imageRelocator: func(controller *gomock.Controller) sheaf.ImageReplacer {
				ir := mocks.NewMockImageReplacer(controller)
				expectBundleManifest("deploy1.yaml", "", ir)
				return ir
			},
			want: "file: deploy1.yaml\n",
		},
		{
			name: "with prefix",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller, []string{"deploy1.yaml"})
			},
			prefix: "prefix",
			imageRelocator: func(controller *gomock.Controller) sheaf.ImageReplacer {
				ir := mocks.NewMockImageReplacer(controller)
				expectBundleManifest("deploy1.yaml", "prefix", ir)
				return ir
			},
			want: "file: deploy1.yaml\n",
		},
		{
			name: "multiple manifests",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller, []string{"deploy1.yaml", "deploy2.yaml"})
			},
			imageRelocator: func(controller *gomock.Controller) sheaf.ImageReplacer {
				ir := mocks.NewMockImageReplacer(controller)
				expectBundleManifest("deploy1.yaml", "", ir)
				expectBundleManifest("deploy2.yaml", "", ir)
				return ir
			},
			want: "file: deploy1.yaml\n---\nfile: deploy2.yaml\n",
		},
		{
			name: "requires bundle factory",
			imageRelocator: func(controller *gomock.Controller) sheaf.ImageReplacer {
				ir := mocks.NewMockImageReplacer(controller)
				return ir
			},
			wantErr: true,
		},
		{
			name: "requires image relocator",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller, []string{"deploy1.yaml", "deploy2.yaml"})
			},
			wantErr: true,
		},
		{
			name: "config get manifests fails",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				bundle := testutil.GenerateBundle(t, controller,
					testutil.BundleGeneratorCreateBundle(
						func(t *testing.T, controller *gomock.Controller, config sheaf.BundleConfig, manifests []sheaf.BundleManifest) *mocks.MockBundle {
							bundle := mocks.NewMockBundle(controller)
							bundle.EXPECT().Config().Return(config)
							bundle.EXPECT().Manifests().Return(nil, fmt.Errorf("error"))
							return bundle
						}))
				return func(string) (sheaf.Bundle, error) {
					return bundle, nil
				}
			},
			imageRelocator: func(controller *gomock.Controller) sheaf.ImageReplacer {
				ir := mocks.NewMockImageReplacer(controller)
				return ir
			},
			wantErr: true,
		},
		{
			name: "manifest list fails",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				bundle := testutil.GenerateBundle(t, controller,
					testutil.BundleGeneratorCreateBundle(
						func(t *testing.T, controller *gomock.Controller, config sheaf.BundleConfig, manifests []sheaf.BundleManifest) *mocks.MockBundle {
							bundle := mocks.NewMockBundle(controller)
							bundle.EXPECT().Config().Return(config)

							ms := mocks.NewMockManifestService(controller)
							bundle.EXPECT().Manifests().Return(ms, nil)

							ms.EXPECT().List().Return(nil, fmt.Errorf("error")).AnyTimes()
							return bundle
						}))
				return func(string) (sheaf.Bundle, error) {
					return bundle, nil
				}
			},
			imageRelocator: func(controller *gomock.Controller) sheaf.ImageReplacer {
				ir := mocks.NewMockImageReplacer(controller)
				return ir
			},
			wantErr: true,
		},
		{
			name: "image relocation fails",
			bundleFactory: func(controller *gomock.Controller) sheaf.BundleFactoryFunc {
				return genBundleFactory(t, controller, []string{"deploy1.yaml", "deploy2.yaml"})
			},
			imageRelocator: func(controller *gomock.Controller) sheaf.ImageReplacer {
				ir := mocks.NewMockImageReplacer(controller)
				ir.EXPECT().Replace(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("error"))
				return ir
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			options := []sheaf.Option{
				sheaf.WithRepositoryPrefix(tc.prefix),
			}

			if tc.bundleFactory != nil {
				options = append(options, sheaf.WithBundleFactory(
					tc.bundleFactory(controller)))
			}

			if tc.imageRelocator != nil {
				options = append(options, sheaf.WithImageReplacer(
					tc.imageRelocator(controller)))
			}

			var buf bytes.Buffer
			options = append(options, sheaf.WithWriter(&buf))

			err := sheaf.ManifestShow(options...)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.want, buf.String())
		})
	}
}
