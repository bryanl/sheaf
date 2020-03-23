/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package remote

import (
	"fmt"
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/stretchr/testify/require"
)

func TestImageReader_Read(t *testing.T) {
	cases := []struct {
		name    string
		wantRef string
		fetcher *fakeFetcher
		options []ImageReaderOption
		wantErr bool
	}{
		{
			name: "in general",
			fetcher: &fakeFetcher{
				image: nil,
				err:   nil,
			},
			options: []ImageReaderOption{},
			wantErr: false,
		},
		{
			name: "insecure registry",
			fetcher: &fakeFetcher{
				image: nil,
				err:   nil,
			},
			options: []ImageReaderOption{
				WithInsecure(true),
			},
			wantErr: false,
		},
		{
			name: "fetcher failed",
			fetcher: &fakeFetcher{
				image: nil,
				err:   fmt.Errorf("error"),
			},
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			options := append(tc.options, func(ir *ImageReader) {
				ir.fetcher = tc.fetcher
			})

			ir := NewImageReader(options...)

			_, err := ir.Read("ref")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, "ref", tc.fetcher.requested.String())
		})
	}
}

type fakeFetcher struct {
	image v1.Image
	err   error

	requested name.Reference
}

var _ Fetcher = &fakeFetcher{}

func (f *fakeFetcher) Fetch(ref name.Reference) (v1.Image, error) {
	f.requested = ref
	return f.image, f.err
}
