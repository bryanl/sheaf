/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/images"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func TestPacker_Pack(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	bundle := genBundle(t, controller)

	layout := mocks.NewMockLayout(controller)
	imageName, err := image.NewName("docker.io/library/image")
	require.NoError(t, err)
	layout.EXPECT().Add(imageName).Return(image.Digest{}, nil)

	layoutFactory := func(root string) (Layout, error) {
		return layout, nil
	}

	a := &fakeArchiver{}

	p := NewPacker(ioutil.Discard,
		PackerLayoutFactory(layoutFactory),
		PackerArchiver(a))

	f, err := ioutil.TempFile("", "archive")
	require.NoError(t, err)

	require.NoError(t, p.Pack(bundle, f))

	require.NoError(t, f.Close())

	// ensure fs config exists
	require.True(t, a.contents.hasKey(sheaf.BundleConfigFilename))

	// ensure manifests exists
	require.True(t, a.contents.hasKey("app", "manifests", "deploy.yaml"))
}

func genBundle(t *testing.T, controller *gomock.Controller) *mocks.MockBundle {
	bundle := mocks.NewMockBundle(controller)

	config := sheaf.BundleConfig{
		Name:          "project",
		Version:       "0.1.0",
		SchemaVersion: "v1alpha1",
	}
	bundle.EXPECT().Config().Return(config).AnyTimes()

	bundleManifests := []sheaf.BundleManifest{
		{
			ID:   "deploy.yaml",
			Data: slurpData(t, filepath.Join("testdata", "manifests", "deploy.yaml")),
		},
	}

	m := mocks.NewMockManifestService(controller)
	m.EXPECT().List().Return(bundleManifests, nil).AnyTimes()

	bundle.EXPECT().Manifests().Return(m, nil).AnyTimes()

	imageList, err := images.New([]string{"image"})
	require.NoError(t, err)
	bundle.EXPECT().Images().Return(imageList, nil).AnyTimes()

	return bundle
}

type fakeArchiver struct {
	contents archiverContents
}

var _ sheaf.Archiver = &fakeArchiver{}

func (f fakeArchiver) Unarchive(src, dest string) error {
	panic("implement me")
}

func (f *fakeArchiver) Archive(src string, w io.Writer) error {
	f.contents = archiverContents{}

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		key := strings.TrimPrefix(path, src+"/")
		f.contents[key] = data

		return nil
	})
}

type archiverContents map[string][]byte

func (ac archiverContents) hasKey(parts ...string) bool {
	_, ok := ac[filepath.Join(parts...)]
	return ok
}
