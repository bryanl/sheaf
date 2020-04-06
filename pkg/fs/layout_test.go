/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_insecureTransport(t *testing.T) {
	tr := newInsecureTransport()

	var scheme string
	tr.roundTripperFunc = func(r *http.Request) (response *http.Response, err error) {
		scheme = r.URL.Scheme

		var b bytes.Buffer

		resp := &http.Response{
			Body: ioutil.NopCloser(&b),
		}

		return resp, nil
	}

	r, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)

	resp, err := tr.RoundTrip(r)
	require.NoError(t, err)

	require.NoError(t, resp.Body.Close())

	require.Equal(t, "http", scheme)
}
