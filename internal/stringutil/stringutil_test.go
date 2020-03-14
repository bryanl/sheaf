/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stringutil

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandom(t *testing.T) {
	actual := Random(6)
	require.Len(t, actual, 6)
}

func TestRandomWithCharset(t *testing.T) {
	actual := RandomWithCharset(6, "a")
	require.Equal(t, "aaaaaa", actual)
}
