// +build !integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/mocks"
)

func Test_writer_Writer(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	a := mocks.NewMockArchiver(controller)
	a.EXPECT().
		Unarchive("archive-path", gomock.Any()).
		Return(nil)

	writerCalled := false
	fw := func(bundlePath, dest string, b bool) error {
		require.Equal(t, "dest", dest)
		writerCalled = true
		return nil
	}

	w := newWriter(
		writerOptionArchiver(a),
		writerOptionFSWriter(fw))

	require.NoError(t, w.Write("archive-path", "dest", false))
	require.True(t, writerCalled)
}
