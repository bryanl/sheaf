/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pivotal/image-relocation/pkg/images"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestBundleConfigWriter_Write(t *testing.T) {
	cases := []struct {
		name     string
		options  func(controller *gomock.Controller) []BundleConfigWriterOption
		config   func(controller *gomock.Controller) *mocks.MockBundleConfig
		expected []byte
		wantErr  bool
	}{
		{
			name: "in general",
			config: func(controller *gomock.Controller) *mocks.MockBundleConfig {
				imageList, err := images.New("foo")
				require.NoError(t, err)

				config := mocks.NewMockBundleConfig(controller)
				config.EXPECT().GetSchemaVersion().Return("v1alpha1")
				config.EXPECT().GetName().Return("test")
				config.EXPECT().GetVersion().Return("0.1.0")
				config.EXPECT().GetImages().Return(&imageList)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{
					{
						APIVersion: "v1",
						Kind:       "Item",
						JSONPath:   "{.}",
					},
				})

				return config
			},
			expected: testutil.Testdata(t, "bundle-config-writer", "bundle.json"),
		},
		{
			name: "without images",
			config: func(controller *gomock.Controller) *mocks.MockBundleConfig {
				config := mocks.NewMockBundleConfig(controller)
				config.EXPECT().GetSchemaVersion().Return("v1alpha1")
				config.EXPECT().GetName().Return("test")
				config.EXPECT().GetVersion().Return("0.1.0")
				config.EXPECT().GetImages().Return(&images.Empty)
				config.EXPECT().GetUserDefinedImages().Return([]sheaf.UserDefinedImage{
					{
						APIVersion: "v1",
						Kind:       "Item",
						JSONPath:   "{.}",
					},
				})

				return config
			},
			expected: testutil.Testdata(t, "bundle-config-writer", "bundle-no-images.json"),
		},
		{
			name: "without user defined images",
			config: func(controller *gomock.Controller) *mocks.MockBundleConfig {
				imageList, err := images.New("foo")
				require.NoError(t, err)

				config := mocks.NewMockBundleConfig(controller)
				config.EXPECT().GetSchemaVersion().Return("v1alpha1")
				config.EXPECT().GetName().Return("test")
				config.EXPECT().GetVersion().Return("0.1.0")
				config.EXPECT().GetImages().Return(&imageList)
				config.EXPECT().GetUserDefinedImages().Return(nil)

				return config
			},
			expected: testutil.Testdata(t, "bundle-config-writer", "bundle-no-udi.json"),
		},
		{
			name: "unable to open destination file",
			config: func(controller *gomock.Controller) *mocks.MockBundleConfig {
				config := mocks.NewMockBundleConfig(controller)
				return config
			},
			options: func(controller *gomock.Controller) []BundleConfigWriterOption {
				return []BundleConfigWriterOption{
					func(bcw *BundleConfigWriter) {
						bcw.openFile = func(s string) (closer io.WriteCloser, err error) {
							return nil, fmt.Errorf("error")
						}
					},
				}
			},
			wantErr: true,
		},
		{
			name: "unable to encode config",
			config: func(controller *gomock.Controller) *mocks.MockBundleConfig {
				config := mocks.NewMockBundleConfig(controller)
				return config
			},
			options: func(controller *gomock.Controller) []BundleConfigWriterOption {
				return []BundleConfigWriterOption{
					func(bcw *BundleConfigWriter) {
						c := mocks.NewMockBundleConfigCodec(controller)
						c.EXPECT().Encode(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error"))
						bcw.codec = c
					},
				}
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			dir, err := ioutil.TempDir("", "sheaf-test")
			require.NoError(t, err)
			defer func() {
				require.NoError(t, os.RemoveAll(dir))
			}()

			bundle := mocks.NewMockBundle(controller)
			bundle.EXPECT().Path().Return(dir)

			var options []BundleConfigWriterOption
			if tc.options != nil {
				options = tc.options(controller)
			}

			bcw := NewBundleConfigWriter(options...)

			err = bcw.Write(bundle, tc.config(controller))
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			got := testutil.SlurpData(t, filepath.Join(dir, "bundle.json"))
			require.Equal(t, string(testutil.NormalizeNewlines(tc.expected)), string(got))
		})
	}
}
