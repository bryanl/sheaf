/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/option"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewSetUserDefinedImage creates a set user defined image command.
func NewSetUserDefinedImage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-udi",
		Short: "Set user defined image in bundle",
		Args:  cobra.NoArgs,
	}

	setupSetUDI(cmd)
	return cmd
}

func setupSetUDI(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ConfigSetUDI, "config-set-udi")
	g.WithBundlePath()
	g.WithUserDefinedImage()
}
