/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package reporter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_std_Header(t *testing.T) {
	header := "header"

	var b bytes.Buffer

	r := newStd(stdOptionOutput(&b))
	r.Header(header)

	actual := b.String()
	expected := "header\n"
	require.Equal(t, expected, actual)
}

func Test_std_Report(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		wanted string
	}{
		{

			name:   "single line",
			in:     "a message",
			wanted: "  a message\n",
		},
		{

			name:   "single line terminated with a new line",
			in:     "a message\n",
			wanted: "  a message\n",
		},
		{

			name:   "multiple lines",
			in:     "line 1\nline 2",
			wanted: "  line 1\n  line 2\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer

			r := newStd(stdOptionOutput(&b))
			r.Report(tc.in)

			actual := b.String()
			require.Equal(t, tc.wanted, actual)

		})
	}
}

func Test_std_Reportf(t *testing.T) {
	cases := []struct {
		name   string
		in     string
		inArgs []interface{}
		wanted string
	}{
		{
			name:   "single line with no args",
			in:     "a message",
			wanted: "  a message\n",
		},
		{
			name:   "single line with trailing newline",
			in:     "a message\n",
			wanted: "  a message\n",
		},
		{
			name:   "single line with args",
			in:     "a message with [%s]",
			inArgs: []interface{}{"item"},
			wanted: "  a message with [item]\n",
		},
		{
			name:   "multiple line with args",
			in:     "line [%d]\nline [%d]",
			inArgs: []interface{}{1, 2},
			wanted: "  line [1]\n  line [2]\n",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var b bytes.Buffer

			r := newStd(stdOptionOutput(&b))
			r.Reportf(tc.in, tc.inArgs...)

			actual := b.String()
			require.Equal(t, tc.wanted, actual)

		})
	}
}
