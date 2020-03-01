/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle_test

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/bundle"
)

func TestGenManifest(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	a := archiver.Default

	mg := bundle.NewManifestGenerator(
		bundle.ManifestGeneratorArchivePath("testdata/gen-manifest.tgz"), // single manifest, image layout omitted
		bundle.ManifestGeneratorPrefix("prefix.com"),
		bundle.ManifestGeneratorArchiver(a))

	b := &strings.Builder{}
	err := mg.Generate(b)
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
