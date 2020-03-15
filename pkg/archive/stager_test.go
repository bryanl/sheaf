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

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/images"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestStager_Stage(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	a := mocks.NewMockArchiver(controller)
	a.EXPECT().Unarchive("archive-uri", gomock.Any()).Return(nil)

	bundle := testutil.GenerateBundle(t, controller)
	bundleFactory := func(string) (sheaf.Bundle, error) {
		return bundle, nil
	}

	imageRelocator := mocks.NewMockImageRelocator(controller)
	imageList, err := images.New([]string{"docker.io/library/image"})
	require.NoError(t, err)
	imageRelocator.EXPECT().
		Relocate(gomock.Any(), "registry-prefix", imageList.Slice(), false).
		Return(nil)

	s := NewStager(
		StagerOptionReporter(reporter.Nop{}),
		StagerOptionArchiver(a),
		StagerOptionBundleFactory(bundleFactory),
		StagerOptionImageRelocator(imageRelocator))

	require.NoError(t, s.Stage("archive-uri", "registry-prefix", false))
}
