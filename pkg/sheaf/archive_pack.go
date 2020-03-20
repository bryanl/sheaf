/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import "fmt"

//go:generate mockgen -destination=../mocks/mock_bundle_packer.go -package mocks github.com/bryanl/sheaf/pkg/sheaf BundlePacker

// ArchivePack packs a bundle.
func ArchivePack(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	bp, err := opts.bundlePacker()
	if err != nil {
		return fmt.Errorf("load bundle packer: %w", err)
	}

	if err := bp.Pack(b, opts.destination, opts.force); err != nil {
		return fmt.Errorf("pack bundle: %w", err)
	}
	return nil
}

// BundlePacker is an interface wrapping the bundle pack command.
type BundlePacker interface {
	// Pack packs a bundle to a destination.
	Pack(b Bundle, dest string, force bool) error
}
