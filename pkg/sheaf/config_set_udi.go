/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
)

// ConfigSetUDI sets a user defined image in the bundle config.
func ConfigSetUDI(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	bcw, err := opts.bundleConfigWriter()
	if err != nil {
		return err
	}

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	udi := opts.userDefinedImage

	config := updateUDI(b.Config(), func(u udiMap) {
		key := UserDefinedImageKey{
			APIVersion: udi.APIVersion,
			Kind:       udi.Kind,
		}
		u[key] = udi
	})

	if err := bcw.Write(b, config); err != nil {
		return fmt.Errorf("write bundle config: %w", err)
	}

	return nil
}
