/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

// ArchiveShowManifests shows manifests in an archive.
func ArchiveShowManifests(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	return withExplodedArchive(opts, func(b Bundle) error {
		optionList = append(optionList, WithBundleFactory(func(rootPath string) (bundle Bundle, err error) {
			return b, nil
		}))

		return ManifestShow(optionList...)
	})
}
