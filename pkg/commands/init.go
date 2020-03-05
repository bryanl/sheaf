/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewInitCommand generates an init command.
func NewInitCommand() *cobra.Command {
	var version string
	var bundlePath string

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize bundle directory",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires path to new fs directory")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			initer := sheaf.NewIniter(
				sheaf.IniterOptionName(args[0]),
				sheaf.IniterOptionVersion(version),
				sheaf.IniterOptionBundlePath(bundlePath))
			return initer.Init()
		},
	}

	cmd.Flags().StringVar(&version, "version", "", "fs version")
	cmd.Flags().StringVar(&bundlePath, "fs-path", "", "fs path")

	return cmd
}
