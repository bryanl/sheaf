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

// NewDeleteUserDefinedImage creates a delete user defined image command.
func NewDeleteUserDefinedImage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-udi",
		Short: "Delete user defined image in bundle",
		Args:  cobra.NoArgs,
	}

	setupDeleteUDI(cmd)
	return cmd
}

func setupDeleteUDI(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ConfigDeleteUDI, "config-delete-udi")
	g.WithBundlePath()
	g.WithUserDefinedImageKey()
}
