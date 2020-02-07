package sheaf

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
)

// CreateConfig is configuration for Create.
type CreateConfig struct {
	Path string
}

// Create creates a bundle.
func Create(config CreateConfig) error {
	// check if directory exists
	fi, err := os.Stat(config.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("bundle directory %q does not exist", config.Path)
		}

		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("%q is not a directory", config.Path)
	}

	bundleConfigFilename := filepath.Join(config.Path, BundleConfigFilename)

	bundleConfig, err := LoadBundleConfig(bundleConfigFilename)
	if err != nil {
		return fmt.Errorf("load bundle config: %w", err)
	}

	spew.Dump(bundleConfig)

	return nil
}
