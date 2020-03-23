/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"bytes"
	"fmt"
)

// ManifestShow shows manifests in a bundle.
func ManifestShow(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	if opts.imageReplacer == nil {
		return fmt.Errorf("image replacer is not configured")
	}

	b, err := opts.bundleFactory(opts.bundlePath)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	config := b.Config()

	ms, err := b.Manifests()
	if err != nil {
		return fmt.Errorf("get manifests service: %w", err)
	}

	manifests, err := ms.List()
	if err != nil {
		return fmt.Errorf("list manifests: %w", err)
	}

	for i, manifest := range manifests {
		data, err := opts.imageReplacer.Replace(manifest, config, opts.repositoryPrefix)
		if err != nil {
			return fmt.Errorf("update manifest %s: %w", manifest.ID, err)
		}

		if i > 0 {
			if _, err := fmt.Fprintln(opts.writer, "---"); err != nil {
				return err
			}
		}

		if _, err := fmt.Fprintln(opts.writer, string(bytes.TrimSpace(data))); err != nil {
			return err
		}
	}

	return nil
}
