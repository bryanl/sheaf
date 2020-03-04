/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/containerd/continuity/fs"

	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

type ManifestService struct {
	manifestsDir string
	reporter     reporter.Reporter
}

var _ sheaf.ManifestService = &ManifestService{}

func NewManifestService(manifestsDir string) (*ManifestService, error) {
	m := ManifestService{
		manifestsDir: manifestsDir,
		reporter:     reporter.Default,
	}

	return &m, nil
}

func (m ManifestService) List() ([]sheaf.BundleManifest, error) {
	entries, err := ioutil.ReadDir(m.manifestsDir)
	if err != nil {
		return nil, fmt.Errorf("read manifests dir %q: %w", m.manifestsDir, err)
	}

	var list []sheaf.BundleManifest

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(m.manifestsDir, entry.Name())

		data, err := ioutil.ReadFile(manifestPath)
		if err != nil {
			return nil, fmt.Errorf("unable to read manifest %q: %w", manifestPath, err)
		}

		bm := sheaf.BundleManifest{
			ID:   manifestPath,
			Data: data,
		}

		list = append(list, bm)
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ID < list[j].ID
	})

	return list, nil
}

func (m ManifestService) Add(manifestPath string) error {
	m.reporter.Header(fmt.Sprintf("Adding manifest %s", manifestPath))

	manifestPath, err := filepath.Abs(manifestPath)
	if err != nil {
		return err
	}

	_, file := filepath.Split(manifestPath)

	if err := os.MkdirAll(m.manifestsDir, 0700); err != nil {
		return err
	}

	dest := filepath.Join(m.manifestsDir, file)
	if _, err := os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			return fs.CopyFile(dest, manifestPath)
		}
		return err
	}

	return fmt.Errorf("%s exists", dest)
}
