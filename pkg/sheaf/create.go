package sheaf

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	// write archive to disk
	outputFile := outputFilename(".", bundleConfig)
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

	if err := Archive(config.Path, f); err != nil {
		return fmt.Errorf("create archive: %w", err)
	}

	return nil
}

func outputFilename(dir string, bundleConfig BundleConfig) string {
	filename := filepath.Join(dir, fmt.Sprintf("%s-%s.tar.gz", bundleConfig.Name, bundleConfig.Version))
	return filename
}
