/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archiver

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/archive"
)

var (
	tarOptions = &archive.TarOptions{
		Compression:      archive.Gzip,
		IncludeFiles:     []string{"."},
		IncludeSourceDir: true,
		NoLchown:         true,
	}
)

// targz creates a gzipped tar archive. Assume src is a directory.
func targz(src string, w io.Writer) error {
	export, err := archive.TarWithOptions(src, tarOptions)
	if err != nil {
		return fmt.Errorf("create tar ball: %w", err)
	}
	defer func() {
		if cErr := export.Close(); cErr != nil {
			log.Printf("close tar ball: %v", err)
		}
	}()

	if _, err := io.Copy(w, export); err != nil {
		return fmt.Errorf("write tar ball: %w", err)
	}

	return nil
}

// untargz unarchives a a gzipped tar archiver.
func untargz(src io.Reader, dst string) error {
	// ungzip
	zr, err := gzip.NewReader(src)
	if err != nil {
		return err
	}
	// untar
	tr := tar.NewReader(zr)

	// uncompress each element
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break // end of archiver
		}
		if err != nil {
			return err
		}

		// validate name against path traversal
		if !validRelPath(header.Name) {
			return fmt.Errorf("tar contained invalid name error %q", header.Name)
		}

		// add dst + re-format slashes according to system
		target := filepath.Join(dst, header.Name)
		// if no join is needed, replace with ToSlash:
		// target = filepath.ToSlash(header.Name)

		// check the type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		// if it's a file create it (with same permission)
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			// copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}
			// manually close here after each file operation; deferring would cause each file close
			// to wait until all operations have completed.

			if err := fileToWrite.Close(); err != nil {
				return err
			}
		}
	}

	return nil
}

// validRelPath checks for path traversal and correct forward slashes
func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}
