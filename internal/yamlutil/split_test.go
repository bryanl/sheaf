/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package yamlutil_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/internal/yamlutil"
	"github.com/stretchr/testify/require"
)

func TestSlicer(t *testing.T) {
	const (
		simple = `---
a: b
`
		simpleWithoutNewline = `---
a: b`
		simpleMadeBare = `a: b
`
		bare = `c: d
`
		bareWithoutNewline = `c: d`
		yamlDirective      = `%YAML 1.2
`
		yamlDirectiveWithoutNewline = `%YAML 1.2`
	)

	cases := []struct {
		name        string
		input       string
		expected    []string
		expectedErr string
	}{
		{
			name:     "empty stream",
			input:    "",
			expected: []string{},
		},
		{
			name:  "simple document",
			input: simple,
			expected: []string{
				simple,
			},
		},
		{ // This test documents the actual behaviour, not the intended behaviour.
			name:  "simple document starting with directive",
			input: yamlDirective + simple,
			expected: []string{
				yamlDirectiveWithoutNewline, // bug?
				simpleMadeBare,
			},
		},
		{ // This test documents the actual behaviour, not the intended behaviour.
			name:  "simple stream",
			input: simple + simple,
			expected: []string{
				simpleWithoutNewline, // where did the newline go?
				simpleMadeBare,
			},
		},
		{ // This test documents the actual behaviour, not the intended behaviour.
			name:  "simple stream with directives",
			input: yamlDirective + simple + yamlDirective + simple,
			expected: []string{
				yamlDirectiveWithoutNewline,                  // bug?
				simpleMadeBare + yamlDirectiveWithoutNewline, // same bug?
				simpleMadeBare,
			},
		},
		{
			name:  "bare document",
			input: bare,
			expected: []string{
				bare,
			},
		},
		{ // This test documents the actual behaviour, not the intended behaviour.
			name:  "stream starting with bare document",
			input: bare + simple,
			expected: []string{
				bareWithoutNewline, // where did the newline go?
				simpleMadeBare,
			},
		},
	}

	for _, tc := range cases {
		actual, err := yamlutil.Split([]byte(tc.input))
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
				return
			}

			a := []string{}
			for _, d := range actual {
				a = append(a, string(d))
			}

			require.Equal(t, tc.expected, a)
		})
	}
}

func TestSlicerWithLargeFile(t *testing.T) {
	expectedFiles := []string{}
	for i := 0; i <= 51; i++ {
		expectedFiles = append(expectedFiles, fmt.Sprintf("file_%d.yaml", i))
	}
	cases := []struct {
		name          string
		inputFile     string
		expectedFiles []string
		expectedErr   string
	}{
		{
			name:          "large file",
			inputFile:     "cert-manager.yaml",
			expectedFiles: expectedFiles,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			contents := testutil.SlurpData(t, filepath.Join("./testdata", tc.inputFile))

			actual, err := yamlutil.Split(contents)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
				return
			}

			a := []string{}
			for _, d := range actual {
				a = append(a, string(d))
			}

			e := []string{}
			for _, f := range tc.expectedFiles {
				d := testutil.SlurpData(t, filepath.Join("./testdata", f))
				e = append(e, string(d))
			}

			require.Equal(t, e, a)
		})
	}
}
