/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/archiver"
	"github.com/bryanl/sheaf/pkg/bundle"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewStageCommand creates a stage command.
func NewStageCommand() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "stage",
		Short: "stage a bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("requires bundle location and registry prefix")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			is := bundle.NewImageRelocator(
				bundle.ImageRelocatorDryRun(dryRun))

			config := sheaf.StageConfig{
				ArchivePath:    args[0],
				RegistryPrefix: args[1],
				BundleFactory:  bundle.DefaultBundleFactory,
				ImageStager:    is,
				Archiver:       archiver.Default,
			}

			return sheaf.Stage(config)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry run")

	return cmd
}
