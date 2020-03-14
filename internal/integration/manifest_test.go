// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"bytes"
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
			withWorkingDirectory(t, func(options wdOptions) {
				b := sheafInit(t, testHarness, "integration", options.dir)

				for _, f := range tc.files {
					err = b.harness.runSheaf(b.dir, defaultSheafRunSettings, "manifest", "add", "-f", f)
					require.NoError(t, err, "unable to add %s", f)
				}

				for _, w := range tc.want {
					cur := filepath.Join(b.dir, w.path)
					checkFileExists(t, cur)
					checkFileMatches(t, cur, w.contents)
				}
			})
		})
	}
}

func Test_sheaf_manifest_show(t *testing.T) {
	td := func(parts ...string) string {
		return testdata(t, append([]string{"manifest", "show"}, parts...)...)
	}

	cases := []struct {
		name      string
		manifests []string
		args      []string
		wanted    []byte
	}{
		{
			name: "show single manifest",
			manifests: []string{
				td("workload1.yaml"),
			},
			wanted: readFile(t, td("single.yaml")),
		},
		{
			name: "show multiple manifests",
			manifests: []string{
				td("workload2.yaml"),
				td("workload1.yaml"),
			},
			wanted: readFile(t, td("multiple.yaml")),
		},
		{
			name: "show manifests with prefix",
			manifests: []string{
				td("workload1.yaml"),
			},
			args: []string{
				"--prefix", "example.com/registry",
			},
			wanted: readFile(t, td("prefix.yaml")),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withWorkingDirectory(t, func(options wdOptions) {
				b := sheafInit(t, testHarness, "integration", options.dir)

				for _, manifest := range tc.manifests {
					_, name := filepath.Split(manifest)
					stageFile(t,
						manifest,
						filepath.Join(b.dir, "app", "manifests", name))
				}

				settings := genSheafRunSettings()
				var actual bytes.Buffer
				settings.Stdout = &actual

				args := append([]string{"manifest", "show"}, tc.args...)

				err := b.harness.runSheaf(b.dir, settings, args...)
				require.NoError(t, err)

				require.Equal(t, string(tc.wanted), actual.String())
			})

		})
	}

}
