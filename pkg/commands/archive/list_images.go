/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/bundle"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewListImages creates a list images command for an archive.
func NewListImages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-images",
		Short: "lists images given an archive path",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires archive path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, err := ioutil.TempDir("", "sheaf")
			if err != nil {
				return err
			}

			defer func() {
				if rErr := os.RemoveAll(dir); rErr != nil {
					log.Printf("remove temporary directory: %v", err)
				}
			}()

			if err := archiver.Default.Unarchive(args[0], dir); err != nil {
				return err
			}

			b, err := bundle.NewBundle(dir)
			if err != nil {
				return err
			}

			config := sheaf.ArchiveListImagesConfig{
				Bundle: b,
			}
			return sheaf.ArchiveListImages(config)
		},
	}

	return cmd
}
