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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// IniterOption is an option for configuring Initer.
type IniterOption func(i Initer) Initer

// IniterOptionBundlePath configures the bundle path for Initer.
func IniterOptionBundlePath(p string) IniterOption {
	return func(i Initer) Initer {
		i.BundlePath = p
		return i
	}
}

// IniterOptionName configure the name for Initer.
func IniterOptionName(name string) IniterOption {
	return func(i Initer) Initer {
		i.Name = name
		return i
	}
}

// IniterOptionVersion configure the version for Initer.
func IniterOptionVersion(name string) IniterOption {
	return func(i Initer) Initer {
		i.Version = name
		return i
	}
}

// Initer initializes a bundle.
type Initer struct {
	BundlePath string
	Name       string
	Version    string
}

// NewIniter creates an instance of Initer.
func NewIniter(options ...IniterOption) *Initer {
	i := Initer{}

	for _, option := range options {
		i = option(i)
	}

	return &i
}

// Init initializes a bundle.
func (i *Initer) Init() error {
	if i.Name == "" {
		return fmt.Errorf("name is blank")
	}

	bc := NewBundleConfig(i.Name, i.Version)

	bundlePath := i.BundlePath
	if bundlePath == "" {
		bundlePath = i.Name
	}

	if err := os.MkdirAll(bundlePath, 0700); err != nil {
		return err
	}

	bundleConfigPath := filepath.Join(bundlePath, BundleConfigFilename)
	f, err := os.Create(bundleConfigPath)
	if err != nil {
		return err
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("close bundle config: %v", err)
		}
	}()

	if err := json.NewEncoder(f).Encode(&bc); err != nil {
		return err
	}

	manifestsPath := filepath.Join(bundlePath, "app", "manifests")
	if err := os.MkdirAll(manifestsPath, 0700); err != nil {
		return err
	}

	return nil
}
