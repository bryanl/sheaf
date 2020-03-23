/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
)

// Init creates a bundle.
func Init(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	if opts.bundleName == "" {
		return fmt.Errorf("bundle name is blank")
	}

	if opts.bundleConfigFactory == nil {
		return fmt.Errorf("bundle config factory is not defined")
	}

	if opts.createBundle == nil {
		return fmt.Errorf("init does not know how to create bundles")
	}

	bc := opts.bundleConfigFactory()
	bc.SetVersion(opts.bundleVersion)
	bc.SetName(opts.bundleName)

	if err := opts.createBundle(bc); err != nil {
		return fmt.Errorf("create bundle: %w", err)
	}

	return nil
}
