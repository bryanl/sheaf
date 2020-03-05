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

// NewAddImageCommand creates an add image command.
func NewAddImageCommand() *cobra.Command {
	var images []string

	cmd := &cobra.Command{
		Use:   "add-image",
		Short: "Add image to bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires fs path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ia, err := sheaf.NewImageAdder(args[0])
			if err != nil {
				return err
			}
			return ia.Add(images)
		},
	}

	cmd.Flags().StringSliceVarP(&images, "image", "i", nil,
		"image to add (can specify multiple times)")

	return cmd
}
