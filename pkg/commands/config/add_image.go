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

// NewAddImageCommand creates an add image command.
func NewAddImageCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add-image",
		Short: "Add image to bundle",
	}

	setupAddImage(cmd)
	return cmd
}

func setupAddImage(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ConfigAddImage, "config-add-image")
	g.WithBundlePath()
	g.WithImages()
}
