/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

//go:generate mockgen -destination=../mocks/mock_image_reader.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ImageReader

// ImageReader is an interface that wraps reading an image from a registry.
type ImageReader interface {
	// Read fetches an image given a reference.
	Read(refStr string) (v1.Image, error)
}
