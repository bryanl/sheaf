// +build !integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"io"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestWriter_Write(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bundlePath := filepath.Join("bundle-path")
	dest := "example.com/bundle"

	a := mocks.NewMockArchiver(controller)
	a.EXPECT().Archive(bundlePath, gomock.Any()).
		DoAndReturn(func(_ string, w io.Writer) error {
			_, err := w.Write([]byte("data"))
			require.NoError(t, err)
			return nil
		})

	imageWritten := false
	var iw sheaf.ImageWriter = func(wantedDest string, i v1.Image, b bool) error {
		require.Equal(t, dest, wantedDest)

		imageWritten = true
		return nil
	}

	w := newWriter(
		writerOptionArchiver(a),
		writerOptionImageWriter(iw),
		writerOptionReporter(reporter.Nop{}))

	err := w.Write(bundlePath, dest, false)
	require.NoError(t, err)

	require.True(t, imageWritten)
}
