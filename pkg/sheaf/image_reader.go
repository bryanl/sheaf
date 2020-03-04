/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ImageReader fetches an image given a reference.
type ImageReader func(refStr string) (v1.Image, error)

// DefaultImageReader fetches an image given a reference from a registry.
func DefaultImageReader(refStr string) (v1.Image, error) {
	ref, err := name.ParseReference(refStr)
	if err != nil {
		return nil, fmt.Errorf("parse remote reference: %w", err)
	}

	return remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
}
