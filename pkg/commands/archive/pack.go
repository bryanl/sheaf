/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/internal/fsutil"
	"github.com/bryanl/sheaf/pkg/fs"
)

// NewPackCommand creates a pack command.
func NewPackCommand() *cobra.Command {
	var bundlePath string
	var force bool

	cmd := &cobra.Command{
		Use:   "pack",
		Short: "Pack a bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("archive pack <bundle directory>")
			} else if len(args) == 1 {
				if !fsutil.IsDirectory(args[0]) {
					return fmt.Errorf("%s is not a directory", args[0])
				}
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			dest := "."
			if len(args) > 0 {
				dest = args[0]
			}

			config := fs.PackConfig{
				BundleURI:     bundlePath,
				BundleFactory: fs.DefaultBundleFactory,
				Packer:        fs.NewPacker(),
				Dest:          dest,
				Force:         force,
			}
			return fs.Pack(config)
		},
	}

	cmd.Flags().StringVar(&bundlePath, "bundle-path", ".", "bundle path")
	cmd.Flags().BoolVar(&force, "force", false, "overwrite existing file")

	return cmd
}
