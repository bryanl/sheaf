/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_fetchURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "guts")
	}))
	defer ts.Close()

	u, err := url.Parse(ts.URL)
	require.NoError(t, err)
	u.Path = "file.txt"

	rc, base, err := fetchURL(u.String())
	require.NoError(t, err)

	got, err := ioutil.ReadAll(rc)
	require.NoError(t, err)
	require.Equal(t, "guts", string(got))

	require.Equal(t, "file.txt", base)

	require.NoError(t, rc.Close())

}
