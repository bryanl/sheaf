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
	"bytes"
	"fmt"

	"github.com/pivotal/image-relocation/pkg/image"
)

func replaceImage(manifest []byte, oldImage image.Name, newImage image.Name) []byte {
	for _, oi := range oldImage.Synonyms() {
		old := fmt.Sprintf("image: %s", oi.String())
		new := fmt.Sprintf("image: %s", newImage.String())
		manifest = bytes.Replace(manifest, []byte(old), []byte(new), -1)

		old = fmt.Sprintf("image: %q", oi.String())
		new = fmt.Sprintf("image: %q", newImage.String())
		manifest = bytes.Replace(manifest, []byte(old), []byte(new), -1)
	}
	return manifest
}
