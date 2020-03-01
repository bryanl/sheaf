/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/bundle"
)

// NewGenManifestCommand creates a gen manifest command.
func NewGenManifestCommand() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "gen-manifest",
		Short: "generate manifest",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires path to archive")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mg := bundle.NewManifestGenerator(
				bundle.ManifestGeneratorArchivePath(args[0]),
				bundle.ManifestGeneratorPrefix(prefix),
				bundle.ManifestGeneratorArchiver(archiver.Default))
			return mg.Generate(os.Stdout)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "registry prefix")

	return cmd
}
