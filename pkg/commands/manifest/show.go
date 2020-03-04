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
		Long: `Print manifest to standard out. With no argument, it will assume you in a bundle directory
or a descendent. Argument can either be a bundle directory or a bundle archive.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p, _ := os.Getwd()
			if len(args) > 0 {
				p = args[0]
			}

			mg := fs.NewManifestShower(
				p,
				fs.ManifestShowerPrefix(prefix),
				fs.ManifestShowerArchiver(archiver.Default))
			return mg.Show(os.Stdout)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "registry prefix")

	return cmd
}
