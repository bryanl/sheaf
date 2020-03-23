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

// ImageReaderOption is a functional option for configuring ImageReader.
// TODO: use remote.Option
type ImageReaderOption func(ir *ImageReader)

// WithInsecure sets an insecure registry.
func WithInsecure(forceInsecure bool) ImageReaderOption {
	return func(ir *ImageReader) {
		ir.forceInsecure = forceInsecure
	}
}

// ImageReader reads images from a remote registry.
type ImageReader struct {
	fetcher         Fetcher
	referenceParser ReferenceParser
	forceInsecure   bool
}

var _ sheaf.ImageReader = &ImageReader{}

// NewImageReader creates an instance of ImageReader.
func NewImageReader(options ...ImageReaderOption) *ImageReader {
	ir := ImageReader{
		fetcher:         &ggcrFetcher{},
		referenceParser: &ggcrReferenceParser{},
	}

	for _, option := range options {
		option(&ir)
	}

	return &ir
}

// Read reads an image from remote.
func (i ImageReader) Read(refStr string) (v1.Image, error) {
	ref, err := i.referenceParser.Parse(refStr, i.forceInsecure)
	if err != nil {
		return nil, fmt.Errorf("parse remote reference: %w", err)
	}

	image, err := i.fetcher.Fetch(ref)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch image %s: %w", refStr, err)
	}

	return image, nil
}

// Fetcher is an interface for fetching a ref.
type Fetcher interface {
	// Fetch fetches a named reference and returns an image.
	Fetch(ref name.Reference) (v1.Image, error)
}

type ggcrFetcher struct{}

var _ Fetcher = &ggcrFetcher{}

func (g ggcrFetcher) Fetch(ref name.Reference) (v1.Image, error) {
	return remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
}

// ReferenceParser parsers a reference.
type ReferenceParser interface {
	// Parse parsers a reference string and returns a named reference.
	Parse(refStr string, forceInsecure bool) (name.Reference, error)
}

type ggcrReferenceParser struct{}

var _ ReferenceParser = &ggcrReferenceParser{}

func (g ggcrReferenceParser) Parse(refStr string, forceInsecure bool) (name.Reference, error) {
	var nameOptions []name.Option
	if forceInsecure {
		nameOptions = append(nameOptions, name.Insecure)
	}

	return name.ParseReference(refStr, nameOptions...)
}
