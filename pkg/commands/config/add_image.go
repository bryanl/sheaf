/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/fs"
)

// NewAddImageCommand creates an add image command.
func NewAddImageCommand() *cobra.Command {
	var images []string
	var bundlePath string

	cmd := &cobra.Command{
		Use:   "add-image",
		Short: "Add image to bundle",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ia, err := fs.NewImageAdder(bundlePath)
			if err != nil {
				return err
			}
			return ia.Add(images...)
		},
	}

	cmd.Flags().StringSliceVarP(&images, "image", "i", nil,
		"image to add (can specify multiple times)")

	cmd.Flags().StringVar(&bundlePath, "bundle-path", ".", "bundle path")

	return cmd
}
