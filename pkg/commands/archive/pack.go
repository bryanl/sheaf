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

// NewPackCommand creates a pack command.
func NewPackCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pack",
		Short: "Pack a bundle",
		Args:  cobra.NoArgs,
	}

	setupPack(cmd)
	return cmd
}

func setupPack(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ArchivePack, "archive-pack")
	g.WithBundlePath()
	g.WithDestination()
	g.WithForce()

}
