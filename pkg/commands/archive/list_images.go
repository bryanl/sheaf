/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/option"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewListImages creates a list images command for an archive.
func NewListImages() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-images",
		Short: "Lists images given an archive path",
		Args:  cobra.NoArgs,
	}

	setupListImages(cmd)
	return cmd
}

func setupListImages(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ArchiveListImages, "archive-list-images")
	g.WithArchive()
	g.WithBundlePath()
}
