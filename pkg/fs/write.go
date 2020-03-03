/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// Write writes a bundle path to a remote location.
func Write(bundlePath string, dest string) error {
	w := newWriter()

	return w.Write(bundlePath, dest)
}

type writerOption func(w writer) writer

func writerOptionArchiver(a sheaf.Archiver) writerOption {
	return func(w writer) writer {
		w.archiver = a
		return w
	}
}

func writerOptionImageWriter(iw sheaf.ImageWriter) writerOption {
	return func(w writer) writer {
		w.imageWriter = iw
		return w
	}
}

func writerOptionReporter(r reporter.Reporter) writerOption {
	return func(w writer) writer {
		w.reporter = r
		return w
	}
}

type writer struct {
	reporter    reporter.Reporter
	archiver    sheaf.Archiver
	imageWriter sheaf.ImageWriter
}

func newWriter(options ...writerOption) *writer {
	w := writer{
		reporter:    reporter.Default,
		archiver:    archiver.Default,
		imageWriter: sheaf.DefaultImageWriter,
	}

	for _, option := range options {
		w = option(w)
	}

	return &w
}

func (w *writer) Write(bundlePath string, dest string) error {
	r := w.reporter

	r.Header("Write configuration to registry")

	archiveBytes, err := w.createArchive(bundlePath)
	if err != nil {
		return err
	}

	var layers []mutate.Addendum

	layer, err := w.createLayer(archiveBytes)
	if err != nil {
		return err
	}

	layers = append(layers, layer)

	r.Report("create image from layers")
	base := empty.Image
	withConfig, err := mutate.Append(base, layers...)
	if err != nil {
		return fmt.Errorf("append data layer to base: %w", err)
	}

	r.Report("setting image configuration")
	cfg, err := withConfig.ConfigFile()
	if err != nil {
		return fmt.Errorf("get config file: %w", err)
	}

	cfg = cfg.DeepCopy()
	cfg.Author = "github.com/bryanl/sheaf"

	r.Report("adding image configuration to image")
	image, err := mutate.ConfigFile(withConfig, cfg)
	if err != nil {
		return fmt.Errorf("mutate config file: %w", err)
	}

	r.Reportf("pushing new image to %s\n", dest)
	return w.imageWriter(dest, image)
}

func (w *writer) createArchive(bundlePath string) ([]byte, error) {
	w.reporter.Reportf("creating archive of configuration and manifests in %s\n", bundlePath)
	a := w.archiver
	var b bytes.Buffer
	if err := a.Archive(bundlePath, &b); err != nil {
		return nil, fmt.Errorf("create archive: %w", err)
	}

	return b.Bytes(), nil
}

func (w *writer) createLayer(b []byte) (mutate.Addendum, error) {
	w.reporter.Report("create layer with archive")
	dataLayer, err := tarball.LayerFromOpener(func() (closer io.ReadCloser, err error) {
		return ioutil.NopCloser(bytes.NewBuffer(b)), nil
	})
	if err != nil {
		return mutate.Addendum{}, fmt.Errorf("create data layer: %w", err)
	}
	return mutate.Addendum{
		Layer: dataLayer,
		History: v1.History{
			Author:    "sheaf",
			CreatedBy: "sheaf <insert command here>",
			Comment:   "experimental",
		},
	}, nil
}
