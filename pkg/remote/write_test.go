/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package remote

import (
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/random"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/reporter"
)

func TestWriter_Write(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	a := mocks.NewMockArchiver(controller)
	a.EXPECT().Unarchive(gomock.Any(), "dest").Return(nil)

	i, err := random.Image(100, 1)
	require.NoError(t, err)
	ir := func(refStr string) (v1.Image, error) {
		return i, nil
	}

	w := newWriter(
		writerOptionReporter(reporter.Nop{}),
		writerOptionArchiver(a),
		writerOptionImageReader(ir))

	require.NoError(t, w.Write("ref", "dest"))
}
