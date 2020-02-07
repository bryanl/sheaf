package sheaf

import (
	"fmt"
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

	if err := bundle.Write(); err != nil {
		return fmt.Errorf("write bundle archive: %w", err)
	}

	return nil
}
