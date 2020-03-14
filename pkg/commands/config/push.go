/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/internal/fsutil"
	"github.com/bryanl/sheaf/pkg/fs"
)

// NewPushCommand create a push command.
func NewPushCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "push",
		Short: "push sheaf bundle config to registry",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("requires path to bundle and destination")
			}

			if !fsutil.IsDirectory(args[0]) {
				return fmt.Errorf("%s is not a directory", args[0])
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			bundlePath := args[0]
			dest := args[1]

			return fs.Write(bundlePath, dest)
		},
	}

	return &cmd
}
