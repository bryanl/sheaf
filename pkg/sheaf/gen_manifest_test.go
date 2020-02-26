/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
