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

	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestImageReplacer_Replace(t *testing.T) {
	cases := []struct {
		name     string
		data     string
		prefix   string
		init     func(t *testing.T, mockBundleConfig *mocks.MockBundleConfig)
		expected string
	}{
		{
			name: "empty prefix",
			data: `spec:
  containers:
  - name: a
    image: b:1`,
			prefix: "",
			init: func(t *testing.T, mockBundleConfig *mocks.MockBundleConfig) {
			},
			expected: `spec:
  containers:
  - name: a
    image: b:1`,
		},
		{
			name: "non-empty prefix",
			data: `spec:
  containers:
  - name: a
    image: b:1`,
			prefix: "example.com/user",
			init: func(t *testing.T, mockBundleConfig *mocks.MockBundleConfig) {
				mockBundleConfig.EXPECT().GetUserDefinedImages().Return(nil)
			},
			expected: `spec:
  containers:
  - name: a
    image: example.com/user/library-b-8b44e6d70542cc94361d2d1db09b8123:1
`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			mockBundleConfig := mocks.NewMockBundleConfig(controller)

			if tc.init != nil {
				tc.init(t, mockBundleConfig)
			}

			imageReplacer := NewImageReplacer()
			actual, err := imageReplacer.Replace(sheaf.BundleManifest{
				ID:   "test",
				Data: []byte(tc.data),
			}, mockBundleConfig, tc.prefix)
			require.NoError(t, err)
			require.Equal(t, tc.expected, string(actual))
		})
	}
}
