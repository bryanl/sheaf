/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package bundle

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
)

//go:generate mockgen -destination=../mocks/mock_layout.go -package mocks github.com/bryanl/sheaf/pkg/bundle Layout

// LayoutFactory creates Layout given a root path.
type LayoutFactory func(root string) (Layout, error)

// DefaultLayoutFactory generates a LayoutFactory.
func DefaultLayoutFactory() LayoutFactory {
	return func(root string) (layout Layout, err error) {
		layoutPath := filepath.Join(root, "artifacts", "layout")
		if _, err := os.Stat(layoutPath); err != nil {
			if os.IsNotExist(err) {
				return ggcr.NewRegistryClient(ggcr.WithTransport(http.DefaultTransport)).
					NewLayout(layoutPath)
			}
			return nil, fmt.Errorf("layout path: %w", err)
		}
		return ggcr.NewRegistryClient(ggcr.WithTransport(http.DefaultTransport)).
			ReadLayout(layoutPath)
	}
}

// Layout manages OCI layouts.
type Layout interface {
	registry.Layout
}
