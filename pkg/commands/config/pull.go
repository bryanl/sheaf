/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/remote"
)

// NewPullCommand create a pull command.
func NewPullCommand() *cobra.Command {
	var forceInsecure bool

	cmd := cobra.Command{
		Use:   "pull",
		Short: "pull sheaf bundle config from registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := args[0]
			dest := args[1]
			return remote.Write(ref, dest, forceInsecure)
		},
	}

	cmd.Flags().BoolVar(&forceInsecure, "insecure-registry", false, "insecure registry")

	return &cmd
}
