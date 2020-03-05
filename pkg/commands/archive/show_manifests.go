/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/fs"
)

// NewShowManifests creates a "show manifests" command.
func NewShowManifests() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "show-manifests",
		Short: "Prints manifest to standard out",
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
