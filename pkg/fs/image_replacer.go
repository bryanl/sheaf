/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"strings"

	"github.com/pivotal/image-relocation/pkg/image"
)

func replaceImage(manifest []byte, imageMap map[image.Name]image.Name) []byte {
	var replacements []string
	for oldImage, newImage := range imageMap {
		for _, oi := range oldImage.Synonyms() {
			replacements = append(replacements, oi.String(), newImage.String())
		}
	}
	return []byte(strings.NewReplacer(replacements...).Replace(string(manifest)))
}
