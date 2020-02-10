package sheaf

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainers(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantErr  bool
		expected []string
	}{
		{
			name:    "invalid file",
			path:    "invalid",
			wantErr: true,
		},
		{
			name: "deployment",
			path: "deployment.yaml",
			expected: []string{
				"nginx:1.7.9",
			},
		},
		{
			name: "pod",
			path: "pod.yaml",
			expected: []string{
				"busybox",
			},
		},
		{
			name:     "service",
			path:     "service.yaml",
			expected: nil,
		},
		{
			name: "multi",
			path: "multi.yaml",
			expected: []string{
				"nginx:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manifestPath := filepath.Join("testdata", tt.path)

			got, err := ContainerImages(manifestPath)
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.expected, got)
		})
	}
}
