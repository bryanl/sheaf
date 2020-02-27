/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"testing"

	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/pathmapping"
	"github.com/stretchr/testify/require"
)

func TestFlattenRepoPathPreserveTagDigest(t *testing.T) {
	tests := []struct {
		name      string
		prefix    string
		imageName string
		want      string
		wantError bool
	}{
		{
			name:      "create new name",
			prefix:    "example.com/project",
			imageName: "gcr.io/project/foo:12345",
			want:      "example.com/project/project-foo-ba7ccd825b7871646277a6b334589b7e:12345",
		},
		{
			name:      "long name",
			prefix:    "example.com/project-" + genString(202),
			imageName: "gcr.io/project/foo:12345",
			want:      "example.com/project-" + genString(202) + "/ba7ccd825b7871646277a6b334589b7e:12345",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			imageName, err := image.NewName(tt.imageName)
			require.NoError(t, err)

			got, err := pathmapping.FlattenRepoPathPreserveTagDigest(tt.prefix, imageName)
			if tt.wantError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got.String())
		})
	}
}

func genString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = 'a'
	}

	return string(b)
}
