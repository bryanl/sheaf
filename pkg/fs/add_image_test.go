// +build !integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/images"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestImageAdder_Add(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bundle := testutil.GenerateBundle(t, controller)

	bf := func(string) (sheaf.Bundle, error) {
		return bundle, nil
	}

	written := false
	writer := func(bundle sheaf.Bundle, config sheaf.BundleConfig) error {
		written = true
		want, err := images.New([]string{"docker.io/library/image"})
		require.NoError(t, err)
		require.Equal(t, &want, config.Images)
		return nil
	}

	option := func() ImageAdderOption {
		return func(ia ImageAdder) ImageAdder {
			ia.bundleFactory = bf
			ia.bundleConfigWriter = writer
			return ia
		}
	}

	ia, err := NewImageAdder(".", option())
	require.NoError(t, err)

	require.NoError(t, ia.Add("image"))
	require.True(t, written)
}
