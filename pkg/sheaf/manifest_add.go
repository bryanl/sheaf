/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import "fmt"

// ManifestAdd adds manifests to a bundle.
func ManifestAdd(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	ms, err := b.Manifests()
	if err != nil {
		return fmt.Errorf("get manifests service: %w", err)
	}

	if err := ms.Add(opts.force, opts.filePaths...); err != nil {
		return fmt.Errorf("unable to add files: %w", err)
	}

	return nil
}
