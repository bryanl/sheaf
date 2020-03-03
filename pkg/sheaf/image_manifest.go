/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// ManifestList is a docker manifest list.
type ManifestList struct {
	SchemaVersion int             `json:"schemaVersion"`
	Manifests     []ImageManifest `json:"manifests"`
}

// ImageManifest represents a docker image.
type ImageManifest struct {
	MediaType   string            `json:"mediaType"`
	Size        int               `json:"size"`
	Digest      string            `json:"digest"`
	Annotations map[string]string `json:"annotations"`
}

// RefName returns the ref name for the image.
func (i *ImageManifest) RefName() string {
	return i.Annotations["org.opencontainers.image.ref.name"]
}

// LoadImageManifestFromIndex loads images from a manifest list.
func LoadImageManifestFromIndex(indexPath string) ([]ImageManifest, error) {
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
