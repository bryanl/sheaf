/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestConfigGet(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bundle := testutil.GenerateBundle(t, controller)

	bcc := mocks.NewMockBundleConfigCodec(controller)
	bcc.EXPECT().
		Encode(gomock.Any(), gomock.Any()).
		DoAndReturn(func(w io.Writer, c sheaf.BundleConfig) error {
			_, err := w.Write([]byte("data"))
			return err
		})

	bundleFactory := func(string) (sheaf.Bundle, error) {
		return bundle, nil
	}

	var buf bytes.Buffer

	options := []sheaf.Option{
		sheaf.WithBundleFactory(bundleFactory),
		sheaf.WithWriter(&buf),
		sheaf.WithBundleConfigCodec(bcc),
	}

	require.NoError(t, sheaf.ConfigGet(options...))

	wanted := "data"
	require.Equal(t, wanted, buf.String())
}
