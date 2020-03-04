/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/fs"
)

// NewAddCommand creates a manifest add command.
func NewAddCommand() *cobra.Command {
	var files []string
	var force bool
	var bundleDir string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "add manifest to bundle",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := fs.NewBundle(bundleDir)
			if err != nil {
				return err
			}

			m, err := b.Manifests()
			if err != nil {
				return err
			}

			for _, file := range files {
				if err := m.Add(file); err != nil {
					return err
				}
			}

			return nil
		},
	}

	// NOTE: what's the worst that can happen?
	cwd, _ := os.Getwd()

	cmd.Flags().StringSliceVarP(&files, "filename", "f", nil,
		"filename to add (can specify multiple times)")
	cmd.Flags().BoolVar(&force, "force", false, "force (overwrite)")
	cmd.Flags().StringVarP(&bundleDir, "bundle-dir", "d", cwd, "bundle directory")

	return cmd
}
