/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

//go:generate mockgen -destination=../mocks/mock_image_writer.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ImageWriter

// ImageWriter is an interface that wraps a Write method.
type ImageWriter interface {
	// Write writes an image to a reference location.
	Write(ref string, image v1.Image) error
	// WriteIndex writes an image index to a reference location.
	WriteIndex(ref string, imageIndex v1.ImageIndex) error
}
