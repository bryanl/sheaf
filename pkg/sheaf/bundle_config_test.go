// +build !integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUserDefinedImage_Validate(t *testing.T) {
	cases := []struct {
		name    string
		in      UserDefinedImage
		wantErr bool
	}{
		{
			name: "valid",
			in: UserDefinedImage{
				APIVersion: "api-version",
				Kind:       "kind",
				JSONPath:   "{.}",
				Type:       SingleResult,
			},
		},
		{
			name: "api version is blank",
			in: UserDefinedImage{
				Kind:     "kind",
				JSONPath: "{.}",
				Type:     SingleResult,
			},
			wantErr: true,
		},
		{
			name: "kind is blank",
			in: UserDefinedImage{
				APIVersion: "api-version",
				JSONPath:   "{.}",
				Type:       SingleResult,
			},
			wantErr: true,
		},
		{
			name: "json path is blank",
			in: UserDefinedImage{
				APIVersion: "api-version",
				Kind:       "kind",
				Type:       SingleResult,
			},
			wantErr: true,
		},
		{
			name: "json path is invalid",
			in: UserDefinedImage{
				APIVersion: "api-version",
				Kind:       "kind",
				JSONPath:   "{.",
				Type:       SingleResult,
			},
			wantErr: true,
		},
		{
			name: "type is blank",
			in: UserDefinedImage{
				APIVersion: "api-version",
				Kind:       "kind",
				JSONPath:   "{.}",
			},
			wantErr: true,
		},
		{
			name: "type is invalid",
			in: UserDefinedImage{
				APIVersion: "api-version",
				Kind:       "kind",
				JSONPath:   "{.}",
				Type:       UserDefinedImageType("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.in.Validate()
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
