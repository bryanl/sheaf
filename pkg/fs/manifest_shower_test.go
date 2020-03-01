/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs_test

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/fs"
)

func TestManifestShower_Show(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	a := archiver.Default

	mg := fs.NewManifestShower(
		fs.ManifestShowerArchivePath("testdata/gen-manifest.tgz"), // single manifest, image layout omitted
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
