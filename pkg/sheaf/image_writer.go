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

// ImageWriter writes an image to a destination..
type ImageWriter func(dest string, image v1.Image, forceInsecure bool) error

// DefaultImageWriter is the the default image writer which writes an image
// to a container registry.
func DefaultImageWriter(dest string, image v1.Image, forceInsecure bool) error {
	var nameOptions []name.Option
	if forceInsecure {
		nameOptions = append(nameOptions, name.Insecure)
	}

	dstRef, err := name.ParseReference(dest, nameOptions...)
	if err != nil {
		return fmt.Errorf("parse remote reference: %w", err)
	}

	return remote.Write(dstRef, image, remote.WithAuthFromKeychain(authn.DefaultKeychain))
}
