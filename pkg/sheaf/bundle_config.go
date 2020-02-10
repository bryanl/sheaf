package sheaf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
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

// Filename returns the bundle archive file name for this BundleConfig.
func (bc *BundleConfig) Filename(dir string) string {
	return filepath.Join(dir, fmt.Sprintf("%s-%s.tgz", bc.Name, bc.Version))
}
