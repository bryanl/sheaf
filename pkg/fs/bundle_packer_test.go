/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/internal/testutil"
	"github.com/bryanl/sheaf/pkg/mocks"
	"github.com/bryanl/sheaf/pkg/reporter"
)

func TestBundlePacker_Pack(t *testing.T) {

	tests := []struct {
		name    string
		bundle  func(controller *gomock.Controller) *mocks.MockBundle
		dest    string
		force   bool
		wantErr bool
	}{
		{
			name: "in general",
			bundle: func(controller *gomock.Controller) *mocks.MockBundle {
				b := testutil.GenerateBundle(t, controller)
				nb := testutil.GenerateBundle(t, controller)

				b.EXPECT().Copy(gomock.Any()).Return(nb, nil)

				return b
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			layout := mocks.NewMockLayout(controller)
			layout.EXPECT().Add(gomock.Any()).Return(image.Digest{}, nil)

			bp := NewBundlePacker(func(bp *BundlePacker) {
				bp.reporter = reporter.Nop{}
				bp.layoutFactory = func(_ string) (Layout, error) {
					return layout, nil
				}
			})

			bundle := test.bundle(controller)
			force := test.force

			tempDir, err := ioutil.TempDir("", "sheaf-test")
			require.NoError(t, err)

			defer func() {
				require.NoError(t, os.RemoveAll(tempDir))
			}()

			dest := filepath.Join(tempDir, test.dest)

			err = bp.Pack(bundle, dest, force)
			if test.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
