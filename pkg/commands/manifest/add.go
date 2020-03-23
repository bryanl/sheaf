/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/option"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewAddCommand creates a manifest add command.
func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "add manifest to bundle",
		Args:  cobra.NoArgs,
	}

	setupAdd(cmd)
	return cmd
}

func setupAdd(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ManifestAdd, "manifest-add")
	g.WithBundlePath()
	g.WithFilePaths()
	g.WithForce()
}
