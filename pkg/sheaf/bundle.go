package sheaf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Bundle represents a bundle
type Bundle struct {
	// Path is the path to the bundle directory.
	Path string
	// Config is the BundleConfig for the bundle.
	Config BundleConfig
}

// LoadBundle loads a bundle.
func LoadBundle(path string) (Bundle, error) {
	// check if directory exists
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Bundle{}, fmt.Errorf("bundle directory %q does not exist", path)
		}

		return Bundle{}, err
	}

	if !fi.IsDir() {
		return Bundle{}, fmt.Errorf("%q is not a directory", path)
	}

	bundleConfigFilename := filepath.Join(path, BundleConfigFilename)

	bundleConfig, err := LoadBundleConfig(bundleConfigFilename)
	if err != nil {
		return Bundle{}, fmt.Errorf("load bundle config: %w", err)
	}

	bundle := Bundle{
		Path:   path,
		Config: bundleConfig,
	}

	return bundle, nil
}

// Bundle writes archive to disk.
func (b *Bundle) Write() error {
	outputFile := b.Config.Filename(".")
	fmt.Println("Creating archive: ", outputFile)
	f, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("unable to write %s: %v", outputFile, err)
		}
	}()

	if err := Archive(b.Path, f); err != nil {
		return fmt.Errorf("create archive: %w", err)
	}

	return nil
}
