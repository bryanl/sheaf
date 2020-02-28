/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
)

// ArchiveListImagesConfig is configuration for ArchiveListImages
type ArchiveListImagesConfig struct {
	Bundle BundleService
}

// ArchiveListImages lists images from a bundle.
func ArchiveListImages(config ArchiveListImagesConfig) error {
	is := config.Bundle.Artifacts().Image()

	list, err := is.List()
	if err != nil {
		return err
	}

	return printImageListAsJSON(list, config.Bundle.Codec())
}

func printImageListAsJSON(list []BundleImage, encoder Encoder) error {
	data, err := encoder.Encode(list)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
