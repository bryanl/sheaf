package sheaf

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFromIndex(t *testing.T) {
	tests := []struct {
		name     string
		index    string
		wantErr  bool
		expected []Image
	}{
		{
			name:  "valid archive",
			index: "index.json",
			expected: []Image{
				{
					MediaType: "application/vnd.docker.distribution.manifest.list.v2+json",
					Size:      1412,
					Digest:    "sha256:ad5552c786f128e389a0263104ae39f3d3c7895579d45ae716f528185b36bc6f",
					Annotations: map[string]string{
						"org.opencontainers.image.ref.name": "docker.io/library/nginx:1.17",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := filepath.Join("testdata", tt.index)

			got, err := LoadFromIndex(index)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}
