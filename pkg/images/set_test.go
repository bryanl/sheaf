/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package images_test

import (
	"testing"

	"github.com/bryanl/sheaf/pkg/images"
	"github.com/stretchr/testify/require"
)

func TestMarshalling(t *testing.T) {
	s, err := images.New([]string{"a", "b"})
	require.NoError(t, err)

	sb, err := s.MarshalJSON()
	require.NoError(t, err)

	var u images.Set
	err = (&u).UnmarshalJSON(sb)
	require.NoError(t, err)

	require.Equal(t, s, u)
}

func TestNullUnmarshalling(t *testing.T) {
	var u images.Set
	err := (&u).UnmarshalJSON([]byte("null"))
	require.NoError(t, err)

	require.Equal(t, images.Empty, u)
}
