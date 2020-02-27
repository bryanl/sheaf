/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewStageCommand creates a stage command.
func NewStageCommand() *cobra.Command {
	var unpackDir string

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
			config := sheaf.StageConfig{
				ArchivePath:    args[0],
				RegistryPrefix: args[1],
				UnpackDir:      unpackDir,
			}

			return sheaf.Stage(config)
		},
	}

	cmd.Flags().StringVar(&unpackDir, "unpack-dir", "", "directory to unpack bundle to")

	return cmd
}
