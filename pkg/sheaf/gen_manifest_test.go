/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenManifest(t *testing.T) {
	mg := NewManifestGenerator(
		ManifestGeneratorArchivePath("testdata/gen-manifest.tgz"), // single manifest, image layout omitted
		ManifestGeneratorPrefix("prefix.com"),
	)

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
