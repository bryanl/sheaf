// +build !integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestUserDefinedImageSetter_Set(t *testing.T) {
	udi1 := sheaf.UserDefinedImage{
		APIVersion: "api-version",
		Kind:       "kind",
		JSONPath:   "{.}",
		Type:       sheaf.MultiResult,
	}
	udi2 := sheaf.UserDefinedImage{
		APIVersion: "api-version2",
		Kind:       "kind",
		JSONPath:   "{.}",
		Type:       sheaf.MultiResult,
	}

	cases := []struct {
		name     string
		existing []sheaf.UserDefinedImage
		item     sheaf.UserDefinedImage
		expected []sheaf.UserDefinedImage
		wantErr  bool
	}{
		{
			name:     "no existing images",
			item:     udi1,
			expected: []sheaf.UserDefinedImage{udi1},
		},
		{
			name:     "with existing image",
			existing: []sheaf.UserDefinedImage{udi1},
			item:     udi1,
			expected: []sheaf.UserDefinedImage{udi1},
		},
		{
			name:     "adding a new image to existing",
			existing: []sheaf.UserDefinedImage{udi1},
			item:     udi2,
			expected: []sheaf.UserDefinedImage{udi1, udi2},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			config := testutil.BundleConfig
			config.UserDefinedImages = tc.existing

			bundle := testutil.GenerateBundle(t, controller,
				testutil.BundleGeneratorConfig(config))

			bf := func(string) (sheaf.Bundle, error) {
				return bundle, nil
			}

			written := false

			writer := func(bundle sheaf.Bundle, config sheaf.BundleConfig) error {
				written = true
				require.Equal(t, tc.expected, config.UserDefinedImages)
				return nil
			}

			option := func() UserDefinedImageSetterOption {
				return func(s UserDefinedImageSetter) UserDefinedImageSetter {
					s.bundleFactory = bf
					s.BundleConfigWriter = writer
					return s
				}
			}

			s := NewUserDefinedImageSetter(option())

			require.NoError(t, s.Set(".", tc.item))
			require.True(t, written)

		})
	}

}
