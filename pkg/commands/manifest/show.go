/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/fs"
)

// NewShowCommand creates a "manifest show" command.
func NewShowCommand() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "show",
		Short: "prints manifest to standard out",
		RunE: func(cmd *cobra.Command, args []string) error {
			mg := fs.NewManifestShower(
				fs.ManifestShowerArchivePath(args[0]),
				fs.ManifestShowerPrefix(prefix),
				fs.ManifestShowerArchiver(archiver.Default))
			return mg.Show(os.Stdout)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "registry prefix")

	return cmd
}
