/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ManifestList is a docker manifest list.
type ManifestList struct {
	SchemaVersion int     `json:"schemaVersion"`
	Manifests     []Image `json:"manifests"`
}

// Image represents a docker image.
type Image struct {
	MediaType   string            `json:"mediaType"`
	Size        int               `json:"size"`
	Digest      string            `json:"digest"`
	Annotations map[string]string `json:"annotations"`
}

// RefName returns the ref name for the image.
func (i *Image) RefName() string {
	return i.Annotations["org.opencontainers.image.ref.name"]
}

// LoadFromIndex loads images from a manifest list.
func LoadFromIndex(indexPath string) ([]Image, error) {
	data, err := ioutil.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("read index: %w", err)
	}

	var list ManifestList
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, fmt.Errorf("decode manifest: %w", err)
	}

	return list.Manifests, nil
}
