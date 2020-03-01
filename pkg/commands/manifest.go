/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/commands/manifest"
)

// NewManifestCommand creates a manifest command.
func NewManifestCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "manifest",
		Short:        "sheaf manifest commands",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		manifest.NewShowCommand())

	return cmd
}
