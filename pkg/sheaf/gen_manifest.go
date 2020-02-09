package sheaf

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// ManifestGeneratorOption is an option for configuring ManifestGenerator.
type ManifestGeneratorOption func(mg ManifestGenerator) ManifestGenerator

// ManifestGeneratorPrefix sets the prefix for relocating images.
func ManifestGeneratorPrefix(p string) ManifestGeneratorOption {
	return func(mg ManifestGenerator) ManifestGenerator {
		mg.Prefix = p
		return mg
	}
}

// ManifestGeneratorArchivePath sets the bundle archive path.
func ManifestGeneratorArchivePath(p string) ManifestGeneratorOption {
	return func(mg ManifestGenerator) ManifestGenerator {
		mg.ArchivePath = p
		return mg
	}
}

// ManifestGenerator generates manifests from a bundle archive.
type ManifestGenerator struct {
	ArchivePath string
	Prefix      string
}

// NewManifestGenerator creates an instance of ManifestGenerator.
func NewManifestGenerator(options ...ManifestGeneratorOption) *ManifestGenerator {
	mg := ManifestGenerator{}

	for _, option := range options {
		mg = option(mg)
	}

	return &mg
}

// Generate generates the manifests contained in a bundle archive to the supplied writer..
func (mg *ManifestGenerator) Generate(w io.Writer) error {
	tmpDir, err := ioutil.TempDir("", "sheaf")
	if err != nil {
		return fmt.Errorf("create temporary directory: %w", err)
	}

	defer func() {
		if rErr := os.RemoveAll(tmpDir); rErr != nil {
			log.Printf("remove generator temporary directory: %v", rErr)
		}
	}()

	unpacker := NewUnpacker(
		UnpackerArchivePath(mg.ArchivePath),
		UnpackerDest(tmpDir))
	if err := unpacker.Unpack(); err != nil {
		return fmt.Errorf("unpack bunle: %w", err)
	}

	manifestsPath := filepath.Join(tmpDir, "app", "manifests")

	entries, err := ioutil.ReadDir(manifestsPath)
	if err != nil {
		return err
	}

	for _, fi := range entries {
		if fi.IsDir() {
			continue
		}

		manifestPath := filepath.Join(manifestsPath, fi.Name())
		data, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintln(w, string(bytes.TrimSpace(data))); err != nil {
			return err
		}
	}

	return nil
}
