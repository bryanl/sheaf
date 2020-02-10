/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

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

	unpacker := NewUnpacker(
		UnpackerArchivePath(config.ArchivePath),
		UnpackerDest(unpackDir))

	if err := unpacker.Unpack(); err != nil {
		return fmt.Errorf("unpack bundle: %w", err)
	}

	bundle, err := OpenBundle(unpackDir)
	if err != nil {
		return fmt.Errorf("open bundle: %w", err)
	}

	defer func() {
		if cErr := bundle.Close(); err != nil {
			log.Printf("close bundle: %v", cErr)
		}
	}()

	layoutPath := filepath.Join(unpackDir, "artifacts", "layout")
	indexPath := filepath.Join(layoutPath, "index.json")

	images, err := LoadFromIndex(indexPath)
	if err != nil {
		return fmt.Errorf("read artifact layout index: %w", err)
	}

	registryClient := ggcr.NewRegistryClient()

	layout, err := registryClient.ReadLayout(layoutPath)
	if err != nil {
		return fmt.Errorf("create registry layout: %w", err)
	}

	for i := range images {
		refName := images[i].RefName()
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

		fmt.Printf("relocating %s to %s\n", imageName.String(), newImageName.String())
		if err := layout.Push(imageDigest, newImageName); err != nil {
			return fmt.Errorf("push %s: %w", newImageName.String(), err)
		}
	}

	return nil
}
