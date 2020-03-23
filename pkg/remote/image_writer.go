/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package remote

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// ImageWriter is a remote image writer.
type ImageWriter struct {
	insecureRegistry bool
}

var _ sheaf.ImageWriter = &ImageWriter{}

// NewImageWriter creates an instance of ImageWriter.
func NewImageWriter(optionList ...Option) *ImageWriter {
	var opts options
	for _, option := range optionList {
		option(&opts)
	}

	iw := &ImageWriter{
		insecureRegistry: opts.insecureRegistry,
	}

	return iw
}

// Write writes an image to a remote registry.
func (i *ImageWriter) Write(ref string, image v1.Image) error {
	var nameOptions []name.Option
	if i.insecureRegistry {
		nameOptions = append(nameOptions, name.Insecure)
	}

	dstRef, err := name.ParseReference(ref, nameOptions...)
	if err != nil {
		return fmt.Errorf("parse remote reference: %w", err)
	}

	return remote.Write(dstRef, image, remote.WithAuthFromKeychain(authn.DefaultKeychain))
}

// WriteIndex writes an index to a remote registry.
func (i *ImageWriter) WriteIndex(ref string, imageIndex v1.ImageIndex) error {
	var nameOptions []name.Option
	if i.insecureRegistry {
		nameOptions = append(nameOptions, name.Insecure)
	}

	dstRef, err := name.ParseReference(ref, nameOptions...)
	if err != nil {
		return fmt.Errorf("parse remote reference: %w", err)
	}

	if err := remote.WriteIndex(dstRef, imageIndex, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		return fmt.Errorf("write index: %w", err)
	}

	return nil
}
