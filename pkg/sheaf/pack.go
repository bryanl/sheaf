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
	"log"
)

// PackConfig is configuration for Pack.
type PackConfig struct {
	Path string
}

// Pack packs a bundle.
func Pack(config PackConfig) error {
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
		if _, err := bundle.Store.Add(ref.String()); err != nil {
			return fmt.Errorf("add %s: %w", ref, err)
		}
	}

	if err := bundle.Write(); err != nil {
		return fmt.Errorf("write bundle archive: %w", err)
	}

	return nil
}
