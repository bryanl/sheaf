package sheaf

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// CreateConfig is configuration for Create.
type CreateConfig struct {
	Path string
}

// Create creates a bundle.
func Create(config CreateConfig) error {
	bundle, err := LoadBundle(config.Path)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	// assume manifests live in `app/manifests`
	manifestsPath := filepath.Join(config.Path, "app", "manifests")
	entries, err := ioutil.ReadDir(manifestsPath)
	if err != nil {
		return fmt.Errorf("read manifests dir %q: %w", manifestsPath, err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(manifestsPath, entry.Name())
		images, err := ContainerImages(manifestPath)
		if err != nil {
			return fmt.Errorf("find container images for %q: %w", manifestPath, err)
		}

		fmt.Printf("Images in %s: [%s]\n", entry.Name(), strings.Join(images, ","))

	}

	if err := bundle.Write(); err != nil {
		return fmt.Errorf("write bundle archive: %w", err)
	}

	return nil
}
