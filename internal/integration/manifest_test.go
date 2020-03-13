// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sheaf_manifest_add(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(filepath.Join("testdata", "manifest", "add", "workload1.yaml"))
		require.NoError(t, err)

		_, err = w.Write(data)
		require.NoError(t, err)
	}))
	defer ts.Close()

	fileUrl, err := url.Parse(ts.URL)
	require.NoError(t, err)

	fileUrl.Path = "workload1.yaml"

	type wanted struct {
		path     string
		contents []byte
	}

	cases := []struct {
		name  string
		files []string
		want  []wanted
	}{
		{
			name: "add file",
			files: []string{
				testdata(t, "manifest", "add", "workload1.yaml"),
			},
			want: []wanted{
				{
					path:     filepath.Join("app", "manifests", "workload1.yaml"),
					contents: readFile(t, testdata(t, "manifest", "add", "workload1.yaml")),
				},
			},
		},
		{
			name: "add multiple file",
			files: []string{
				testdata(t, "manifest", "add", "workload1.yaml"),
				testdata(t, "manifest", "add", "workload2.yaml"),
			},
			want: []wanted{
				{
					path:     filepath.Join("app", "manifests", "workload1.yaml"),
					contents: readFile(t, testdata(t, "manifest", "add", "workload1.yaml")),
				},
				{
					path:     filepath.Join("app", "manifests", "workload2.yaml"),
					contents: readFile(t, testdata(t, "manifest", "add", "workload2.yaml")),
				},
			},
		},
		{
			name: "add directory",
			files: []string{
				testdata(t, "manifest", "add"),
			},
			want: []wanted{
				{
					path:     filepath.Join("app", "manifests", "workload1.yaml"),
					contents: readFile(t, testdata(t, "manifest", "add", "workload1.yaml")),
				},
			},
		},
		{
			name: "add URL",
			files: []string{
				fileUrl.String(),
			},
			want: []wanted{
				{
					path:     filepath.Join("app", "manifests", "workload1.yaml"),
					contents: readFile(t, testdata(t, "manifest", "add", "workload1.yaml")),
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withWorkingDirectory(t, func(wd string) {
				err := testHarness.runSheaf(wd, "init", "integration")
				require.NoError(t, err)

				bundleDir := filepath.Join(wd, "integration")

				for _, f := range tc.files {
					err = testHarness.runSheaf(bundleDir, "manifest", "add", "-f", f)
					require.NoError(t, err, "unable to add %s", f)
				}

				for _, w := range tc.want {
					cur := filepath.Join(bundleDir, w.path)
					checkFileExists(t, cur)
					checkFileMatches(t, cur, w.contents)
				}
			})
		})
	}
}
