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

// NewShowCommand creates a "manifest show" command.
func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "prints manifest to standard out",
		Long: `Print manifest to standard out. With no argument, it will assume you in a bundle directory
or a descendent. Argument can either be a bundle directory or a bundle archive.`,
		Args: cobra.NoArgs,
	}

	setupShow(cmd)
	return cmd
}

func setupShow(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ManifestShow, "manifest-shwo")
	g.WithBundlePath()
	g.WithPrefix()
}
