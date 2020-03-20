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

// NewGetCommand creates `config get` command.
func NewGetCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get bundle configuration",
		Args:  cobra.NoArgs,
	}

	setupGet(&cmd)

	return &cmd
}

func setupGet(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ConfigGet, "config-get")
	g.WithBundlePath()
}
