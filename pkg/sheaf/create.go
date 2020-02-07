package sheaf

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
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

	spew.Dump(images)

	if err := bundle.Write(); err != nil {
		return fmt.Errorf("write bundle archive: %w", err)
	}

	return nil
}
