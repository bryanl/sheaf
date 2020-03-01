/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/mocks"
)

func TestArtifactsService_Index(t *testing.T) {
	dir, err := ioutil.TempDir("", "test")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.RemoveAll(dir))
	}()

	layoutDir := filepath.Join(dir, "artifacts", "layout")
	require.NoError(t, os.MkdirAll(layoutDir, 0700))

	wanted := stageFile(t, "index.json", filepath.Join(layoutDir, "index.json"))

	controller := gomock.NewController(t)
	defer controller.Finish()

	bundleService := mocks.NewMockBundle(controller)
	bundleService.EXPECT().Path().Return(dir)

	as := NewArtifactsService(bundleService)

	actual, err := as.Index()
	require.NoError(t, err)
	require.Equal(t, wanted, actual)
}
