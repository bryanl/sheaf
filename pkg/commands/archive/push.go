/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/archive"
)

// NewPushCommand create a push command.
func NewPushCommand() *cobra.Command {
	var forceInsecure bool

	cmd := cobra.Command{
		Use: "push",
		// TODO: support relocating images
		Short: "Push sheaf archive to registry",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("requires path to bundle and destination")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			bundlePath := args[0]
			dest := args[1]

			return archive.Write(bundlePath, dest, forceInsecure)
		},
	}

	cmd.Flags().BoolVar(&forceInsecure, "insecure-registry", false, "insecure registry")

	return &cmd
}
