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
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestUserDefinedImageDeleter_Delete(t *testing.T) {
	udi1 := sheaf.UserDefinedImage{
		APIVersion: "api-version",
		Kind:       "kind",
		JSONPath:   "{.}",
		Type:       sheaf.MultiResult,
	}

	controller := gomock.NewController(t)
	defer controller.Finish()

	config := testutil.BundleConfig
	config.UserDefinedImages = []sheaf.UserDefinedImage{udi1}

	bundle := testutil.GenerateBundle(t, controller,
		testutil.BundleGeneratorConfig(config))

	bf := func(string) (sheaf.Bundle, error) {
		return bundle, nil
	}

	written := false
	writer := func(bundle sheaf.Bundle, config sheaf.BundleConfig) error {
		written = true
		require.Empty(t, config.UserDefinedImages)
		return nil
	}

	option := func() UserDefinedImageDeleterOption {
		return func(d UserDefinedImageDeleter) UserDefinedImageDeleter {
			d.bundleFactory = bf
			d.bundleConfigWriter = writer
			return d
		}
	}

	d := NewUserDefinedImageDeleter(option())

	key := sheaf.UserDefinedImageKey{
		APIVersion: udi1.APIVersion,
		Kind:       udi1.Kind,
	}

	require.NoError(t, d.Delete(".", key))
	require.True(t, written)
}
