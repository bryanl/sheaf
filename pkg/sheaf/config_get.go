/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
)

// ConfigGet shows the bundle configuration.
func ConfigGet(options ...Option) error {
	opts := makeDefaultOptions(options...)

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	if err := opts.bundleConfigCodec.Encode(opts.writer, b.Config()); err != nil {
		return fmt.Errorf("generate bundle config: %w", err)
	}

	return nil
}
