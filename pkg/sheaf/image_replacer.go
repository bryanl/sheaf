/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"strings"

	"github.com/pivotal/image-relocation/pkg/image"
)

func replaceImage(manifest []byte, imageMap map[image.Name]image.Name) []byte {
	replacements := []string{}
	for oldImage, newImage := range imageMap {
		for _, oi := range oldImage.Synonyms() {
			replacements = append(replacements, oi.String(), newImage.String())
		}
	}
	return []byte(strings.NewReplacer(replacements...).Replace(string(manifest)))
}
