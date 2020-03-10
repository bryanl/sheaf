/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/commands/archive"
)

// NewArchiveCommand creates an archive command.
func NewArchiveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "archive",
		Short:        "Perform actions on an archive",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		archive.NewListImages(),
		archive.NewPackCommand(),
		archive.NewPushCommand(),
		archive.NewStageCommand(),
		archive.NewShowManifests())

	return cmd
}
