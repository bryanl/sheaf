/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"

	"github.com/pivotal/image-relocation/pkg/images"
)

// ConfigAddImage adds an image to a bundle configuration.
func ConfigAddImage(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	bcw, err := opts.bundleConfigWriter()
	if err != nil {
		return err
	}

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	bc := b.Config()

	im, err := images.New(opts.images...)
	if err != nil {
		return err
	}

	cur := images.Empty
	if imageSet := bc.GetImages(); imageSet != nil {
		cur = *imageSet
	}

	n := cur.Union(im)
	bc.SetImages(&n)

	if err := bcw.Write(b, bc); err != nil {
		return fmt.Errorf("write bundle config: %w", err)
	}

	return nil
}
