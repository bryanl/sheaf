/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
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

// IniterOptionBundlePath configures the fs path for Initer.
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

// Initer initializes a fs.
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

// Init initializes a fs.
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
			log.Printf("close fs config: %v", err)
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
