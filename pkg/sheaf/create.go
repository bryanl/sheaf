package sheaf

import (
	"fmt"
	"log"
)

// CreateConfig is configuration for Create.
type CreateConfig struct {
	Path string
}

// Create creates a bundle.
func Create(config CreateConfig) error {
	bundle, err := OpenBundle(config.Path)
	if err != nil {
		return fmt.Errorf("load bundle: %w", err)
	}

	defer func() {
		if cErr := bundle.Close(); cErr != nil {
			log.Printf("unable to close bundle: %v", err)
		}
	}()

	images, err := bundle.Images()
	if err != nil {
		return fmt.Errorf("collect images from manifest: %w", err)
	}

	for _, ref := range images {
		fmt.Printf("Adding %s to bundle\n", ref)
		if _, err := bundle.Store.Add(ref); err != nil {
			return fmt.Errorf("add %s: %w", ref, err)
		}
	}

	if err := bundle.Write(); err != nil {
		return fmt.Errorf("write bundle archive: %w", err)
	}

	return nil
}
