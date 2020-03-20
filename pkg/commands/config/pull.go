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

// NewPullCommand create a pull command.
func NewPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull",
		Short: "pull sheaf bundle config from registry",
		Args:  cobra.NoArgs,
	}

	setupPull(cmd)
	return cmd
}

func setupPull(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ConfigPull, "config-pull")
	g.WithInsecureRegistry()
	g.WithReference()
	g.WithDestination()
}
