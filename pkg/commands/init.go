/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/option"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewInitCommand generates an init command.
func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize bundle directory",
		Args:  cobra.NoArgs,
	}

	setupInit(cmd)

	return cmd
}

func setupInit(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.Init, "init")
	g.WithBundleName()
	g.WithInitBundlePath()
	g.WithBundleVersion()
	g.WithBundleConfigFactory()
}
