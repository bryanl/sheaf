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
	"path/filepath"

	"github.com/bryanl/sheaf/pkg/images"
)

const (
	// BundleConfigFilename is the filename for a bundle config.
	BundleConfigFilename = "bundle.json"

	bundleConfigDefaultVersion = "0.1.0"
)

// BundleConfig is a bundle configuration.
type BundleConfig struct {
	// Name is the name of the bundle.
	Name string `json:"name"`
	// Version is the version of the bundle.
	Version string `json:"version"`
	// SchemaVersion is the version of the schema this bundle uses.
	SchemaVersion string `json:"schemaVersion"`
	// Images is a set of images required by the bundle.
	Images images.Set `json:"images"`
}

// NewBundleConfig creates a BundleConfig.
func NewBundleConfig(name, version string) BundleConfig {
	if version == "" {
		version = bundleConfigDefaultVersion
	}

	return BundleConfig{
		Name:          name,
		Version:       version,
		SchemaVersion: "v1alpha1",
	}
}

// LoadBundleConfig loads a BundleConfig from a file.
func LoadBundleConfig(filename string) (BundleConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return BundleConfig{}, fmt.Errorf("read %q: %w", filename, err)
	}

	var bc BundleConfig
	if err := json.Unmarshal(data, &bc); err != nil {
		return BundleConfig{}, err
	}

	return bc, nil
}

// StoreBundleConfig saves a BundleConfig to a file, destructively.
func StoreBundleConfig(bc BundleConfig, filename string) error {
	jbc, err := json.Marshal(bc)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, jbc, 0644)
}

// Filename returns the bundle archive file name for this BundleConfig.
func (bc *BundleConfig) Filename(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s-%s.tgz", bc.Name, bc.Version))
}
