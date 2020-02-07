package sheaf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	// BundleConfigFilename is the filename for a bundle config.
	BundleConfigFilename = "bundle.json"
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
