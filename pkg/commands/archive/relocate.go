/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/archive"
	"github.com/bryanl/sheaf/pkg/fs"
)

// NewStageCommand creates a stage command.
func NewStageCommand() *cobra.Command {
	var dryRun bool
	var forceInsecure bool

	cmd := &cobra.Command{
		Use:   "relocate",
		Short: "Relocate images in archive to new registry",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("requires fs location and registry prefix")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var bfOptions []fs.LayoutOptionFunc
			if forceInsecure {
				bfOptions = append(bfOptions, fs.DefaultLayoutFactoryInsecureSkipVerify())
			}

			bf := fs.DefaultLayoutFactory(bfOptions...)

			relocator := fs.NewImageRelocator(
				fs.ImageRelocatorDryRun(dryRun),
				fs.ImageRelocatorLayoutFactory(bf))

			stager := archive.NewStager(
				archive.StagerOptionImageRelocator(relocator))

			return stager.Stage(args[0], args[1], forceInsecure)
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "dry run")
	cmd.Flags().BoolVar(&forceInsecure, "insecure-registry", false, "insecure registry")

	return cmd
}
