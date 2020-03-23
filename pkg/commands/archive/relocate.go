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

// NewStageCommand creates a stage command.
func NewStageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "relocate",
		Short: "Relocate images in archive to new registry",
		Args:  cobra.NoArgs,
	}

	setupRelocate(cmd)
	return cmd
}

func setupRelocate(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ArchiveRelocate, "archive-relocate")
	g.WithBundlePath()
	g.WithArchive()
	g.WithInsecureRegistry()
	g.WithPrefix()
	g.WithDryRun()

}
