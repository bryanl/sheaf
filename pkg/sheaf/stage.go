package sheaf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
)

// StageConfig is configuration for Stage.
type StageConfig struct {
	ArchivePath    string
	RegistryPrefix string
	UnpackDir      string
}

// Stage stages an archive in a registry.
func Stage(config StageConfig) error {
	unpackDir := config.UnpackDir
	if unpackDir == "" {
		fmt.Println("using temporary directory to unpack")
		tmpDir, err := ioutil.TempDir("", "sheaf")
		if err != nil {
			return fmt.Errorf("create temp dir: %w", err)
		}
		unpackDir = tmpDir

		defer func() {
			if rErr := os.RemoveAll(unpackDir); rErr != nil {
				log.Printf("remove temporary bundle path %q: %v", unpackDir, rErr)
			}
		}()
	} else {
		fmt.Printf("using unpack directory %s\n", unpackDir)
		if err := os.MkdirAll(unpackDir, 0700); err != nil {
			return fmt.Errorf("unable to create unpack dir: %w", err)
		}
	}

	// TODO: split this up. unarchive and then open bundle
	bundle, err := ImportBundle(config.ArchivePath, unpackDir)
	if err != nil {
		return fmt.Errorf("import bundle from %q: %w", config.ArchivePath, err)
	}

	defer func() {
		if cErr := bundle.Close(); err != nil {
			log.Printf("close bundle: %v", cErr)
		}
	}()

	layoutPath := filepath.Join(bundle.Path, "artifacts", "layout")
	indexPath := filepath.Join(layoutPath, "index.json")

	images, err := LoadFromIndex(indexPath)
	if err != nil {
		return fmt.Errorf("read artifact layout index: %w", err)
	}

	spew.Dump(images)

	registryClient := ggcr.NewRegistryClient()

	fmt.Println("layoutPath", layoutPath)

	layout, err := registryClient.ReadLayout(layoutPath)
	if err != nil {
		return fmt.Errorf("create registry layout: %w", err)
	}

	for i := range images {
		refName := images[i].Annotations["org.opencontainers.image.ref.name"]
		imageName, err := image.NewName(refName)
		if err != nil {
			return fmt.Errorf("create image name for ref %q: %w", refName, err)
		}
		imageDigest, err := layout.Find(imageName)
		if err != nil {
			return fmt.Errorf("find image digest for ref %q: %w", refName, err)
		}

		newImageName, err := FlattenRepoPathPreserveTagDigest(config.RegistryPrefix, imageName)
		if err != nil {
			return fmt.Errorf("create relocated image name: %w", err)
		}

		spew.Dump(imageDigest, newImageName)

		if err := layout.Push(imageDigest, newImageName); err != nil {
			return fmt.Errorf("push %s: %w", newImageName.String(), err)
		}
	}

	return nil
}
