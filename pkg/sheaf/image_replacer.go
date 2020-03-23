/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

//go:generate mockgen -destination=../mocks/mock_image_replacer.go -package mocks github.com/bryanl/sheaf/pkg/sheaf ImageReplacer

// ImageReplacer is an interface that wraps images replacing.
type ImageReplacer interface {
	// Replace replaces images in a bundle manifest with prefixed options.
	Replace(manifest BundleManifest, config BundleConfig, prefix string) ([]byte, error)
}
