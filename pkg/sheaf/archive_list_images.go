/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// ArchiveListImages lists images in an archive.
func ArchiveListImages(optionList ...Option) error {
	opts := makeDefaultOptions(optionList...)

	return withExplodedArchive(opts, func(b Bundle) error {
		is := b.Artifacts().Image()
		list, err := is.List()
		if err != nil {
			return err
		}

		data, err := opts.codec.Encode(list)
		if err != nil {
			return fmt.Errorf("encode image list: %w", err)
		}

		fmt.Fprintln(opts.writer, string(data))

		return nil
	})
}

func withExplodedArchive(opts options, f func(b Bundle) error) error {
	if opts.archive == "" {
		return fmt.Errorf("archive path is required")
	}

	b, err := explodeArchive(opts.archive, opts.archiver, opts.bundleFactory)
	if err != nil {
		return fmt.Errorf("explode archive %s: %w", opts.archive, err)
	}

	defer func() {
		if rErr := os.RemoveAll(b.Path()); rErr != nil {
			log.Printf("remove temporary directory: %v", err)
		}
	}()

	return f(b)
}

func explodeArchive(archive string, archiver Archiver, bundleFactory BundleFactoryFunc) (Bundle, error) {
	dir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return nil, err
	}

	if err := archiver.UnarchivePath(archive, dir); err != nil {
		return nil, fmt.Errorf("unable to unarchive %s: %w", archive, err)
	}

	b, err := bundleFactory(dir)
	if err != nil {
		return nil, fmt.Errorf("load bundle: %w", err)
	}

	return b, nil
}
