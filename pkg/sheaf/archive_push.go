/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

// ArchivePush push an archive's bundle to a registry.
func ArchivePush(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	return withExplodedArchive(opts, func(b Bundle) error {
		optionList = append(optionList, WithBundleFactory(func(rootPath string) (bundle Bundle, err error) {
			return b, nil
		}))

		return ConfigPush(optionList...)
	})
}
