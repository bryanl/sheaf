package sheaf

import "github.com/davecgh/go-spew/spew"

// CreateConfig is configuration for Create.
type CreateConfig struct {
	Path string
}

// Create creates a bundle.
func Create(config CreateConfig) error {
	spew.Dump(config)

	return nil
}
