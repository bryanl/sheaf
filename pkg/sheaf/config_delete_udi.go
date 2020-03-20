/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import "fmt"

// ConfigDeleteUDI deletes a user defined image from a bundle configuration.
func ConfigDeleteUDI(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	bcw, err := opts.bundleConfigWriter()
	if err != nil {
		return err
	}

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	config := updateUDI(b.Config(), func(u udiMap) {
		delete(u, opts.userDefinedImageKey)
	})

	if err := bcw.Write(b, config); err != nil {
		return fmt.Errorf("write bundle config: %w", err)
	}

	return nil
}
