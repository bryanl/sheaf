/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package remote

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/reporter"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// Write writes a remote reference to a local destination.
func Write(ref, dest string, forceInsecure bool) error {
	w := newWriter()

	return w.Write(ref, dest, forceInsecure)
}

type writerOption func(w writer) writer

func writerOptionReporter(r reporter.Reporter) writerOption {
	return func(w writer) writer {
		w.reporter = r
		return w
	}
}

func writerOptionArchiver(a sheaf.Archiver) writerOption {
	return func(w writer) writer {
		w.archiver = a
		return w
	}
}

func writerOptionImageReader(ir sheaf.ImageReader) writerOption {
	return func(w writer) writer {
		w.imageReader = ir
		return w
	}
}

type writer struct {
	reporter    reporter.Reporter
	archiver    sheaf.Archiver
	imageReader sheaf.ImageReader
}

func newWriter(options ...writerOption) *writer {
	w := writer{
		reporter:    reporter.Default,
		archiver:    archiver.Default,
		imageReader: sheaf.DefaultImageReader,
	}

	for _, option := range options {
		w = option(w)
	}

	return &w
}

func (w *writer) Write(refStr, dest string, forceInsecure bool) error {
	_, err := os.Stat(dest)
	if err == nil {
		return fmt.Errorf("destination %s already exists", dest)
	} else if !os.IsNotExist(err) {
		return err
	}

	w.reporter.Headerf("Validating reference")

	w.reporter.Reportf("Pulling image %s from registry", refStr)

	image, err := w.imageReader(refStr, forceInsecure)
	if err != nil {
		return err
	}

	w.reporter.Reportf("Extracting layers from image")
	layers, err := image.Layers()
	if err != nil {
		return err
	}

	if len(layers) != 1 {
		return fmt.Errorf("invalid image format: expected 1 layer; got %d layers", len(layers))
	}

	layer := layers[0]

	rc, err := layer.Compressed()
	if err != nil {
		return err
	}

	defer func() {
		if cErr := rc.Close(); cErr != nil {
			log.Printf("close image layer: %v", err)
		}
	}()

	f, err := ioutil.TempFile("", "image")
	if err != nil {
		return err
	}

	if _, err := io.Copy(f, rc); err != nil {
		if cErr := f.Close(); cErr != nil {
			log.Printf("close temporary file: %v", err)
		}

		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	w.reporter.Reportf("Decompressing image")
	if err := w.archiver.Unarchive(f.Name(), dest); err != nil {
		return err
	}

	return nil
}
