/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import v1 "github.com/google/go-containerregistry/pkg/v1"

//go:generate mockgen -destination=../mocks/mock_bundle_imager.go -package mocks github.com/bryanl/sheaf/pkg/sheaf BundleImager

// BundleImager is an interface that wraps the create image from bundle functionality.
type BundleImager interface {
	// CreateImage creates an image from a bundle.
	CreateImage(b Bundle) (v1.Image, error)
}
