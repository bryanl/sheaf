/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs_test

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/fs"
)

func TestManifestShower_Show_archive(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	a := archiver.Default

	mg := fs.NewManifestShower(
		filepath.Join("testdata", "gen-manifest.tgz"),
		fs.ManifestShowerPrefix("prefix.com"),
		fs.ManifestShowerArchiver(a))

	b := &strings.Builder{}
	err := mg.Show(b)
	require.NoError(t, err)

	require.Equal(t, `apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-deployment
spec:
  template:
     spec:
      containers:
      - name: sample-container
        image: prefix.com/library-nginx-dba37485fee3d4d76d5d82609cc9bccb
        ports:
        - containerPort: 80
`, b.String())
}

func TestManifestShower_Show_fs(t *testing.T) {
	expectedManifest := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-deployment
spec:
  template:
     spec:
      containers:
      - name: sample-container
        image: prefix.com/library-nginx-dba37485fee3d4d76d5d82609cc9bccb
        ports:
        - containerPort: 80
`

	cases := []struct {
		name string
		dir  string
	}{
		{
			name: "in root of bundle",
		},
		{
			name: "in nested directory of bundle",
			dir:  "app",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			a := archiver.Default

			tempDir, err := ioutil.TempDir("", "sheaf-test")
			require.NoError(t, err)

			require.NoError(t, a.Unarchive(filepath.Join("testdata", "gen-manifest.tgz"), tempDir))

			mg := fs.NewManifestShower(
				filepath.Join(tempDir, tc.dir),
				fs.ManifestShowerPrefix("prefix.com"),
				fs.ManifestShowerArchiver(a))

			b := &strings.Builder{}
			err = mg.Show(b)
			require.NoError(t, err)

			require.Equal(t, expectedManifest, b.String())

		})
	}
}
