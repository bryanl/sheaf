/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/pivotal/image-relocation/pkg/registry/ggcr"
	"github.com/pivotal/image-relocation/pkg/transport"
)

//go:generate mockgen -destination=../mocks/mock_layout.go -package mocks github.com/bryanl/sheaf/pkg/fs Layout

// LayoutFactory creates Layout given a root path.
type LayoutFactory func(root string) (Layout, error)

// LayoutOptionFunc is a functional option for configuring DefaultLayoutFactory.
type LayoutOptionFunc func(options LayoutOptions) LayoutOptions

// LayoutOptions are options for DefaultLayoutFactory.
type LayoutOptions struct {
	insecureSkipVerify bool
	certs              []string
}

// DefaultLayoutFactoryInsecureSkipVerify configures support for insecure registries.
func DefaultLayoutFactoryInsecureSkipVerify() LayoutOptionFunc {
	return func(options LayoutOptions) LayoutOptions {
		options.insecureSkipVerify = true
		return options
	}
}

// DefaultLayoutFactory generates a LayoutFactory.
func DefaultLayoutFactory(options ...LayoutOptionFunc) LayoutFactory {
	var lo LayoutOptions
	for _, option := range options {
		lo = option(lo)
	}

	return func(root string) (layout Layout, err error) {
		var t http.RoundTripper

		if lo.insecureSkipVerify {
			t = newInsecureTransport()
		} else {
			nt, err := transport.NewHttpTransport(lo.certs, lo.insecureSkipVerify)
			if err != nil {
				return nil, fmt.Errorf("create http transport: %w", err)
			}

			t = nt
		}

		layoutPath := filepath.Join(root, "artifacts", "layout")
		if _, err := os.Stat(layoutPath); err != nil {
			if os.IsNotExist(err) {
				return ggcr.NewRegistryClient(ggcr.WithTransport(t)).
					NewLayout(layoutPath)
			}
			return nil, fmt.Errorf("layout path: %w", err)
		}
		return ggcr.NewRegistryClient(ggcr.WithTransport(t)).
			ReadLayout(layoutPath)
	}
}

// Layout manages OCI layouts.
type Layout interface {
	registry.Layout
}

type insecureTransport struct {
	roundTripperFunc func(*http.Request) (*http.Response, error)
}

var _ http.RoundTripper = &insecureTransport{}

func newInsecureTransport() *insecureTransport {
	return &insecureTransport{
		roundTripperFunc: http.DefaultTransport.RoundTrip,
	}
}

func (i insecureTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.URL.Scheme = "http"
	return i.roundTripperFunc(r)
}
