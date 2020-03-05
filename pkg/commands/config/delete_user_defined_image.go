/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/fs"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewDeleteUserDefinedImage creates a delete user defined image command.
func NewDeleteUserDefinedImage() *cobra.Command {
	var bundlePath string
	var apiVersion string
	var kind string

	cmd := &cobra.Command{
		Use:   "delete-udi",
		Short: "Delete user defined image in bundle",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			d := fs.NewUserDefinedImageDeleter()

			key := sheaf.UserDefinedImageKey{
				APIVersion: apiVersion,
				Kind:       kind,
			}

			return d.Delete(bundlePath, key)
		},
	}

	cmd.Flags().StringVar(&apiVersion, "api-version", "", "api version")
	cmd.Flags().StringVar(&kind, "kind", "", "kind")

	cmd.Flags().StringVar(&bundlePath, "bundle-path", ".", "bundle path")

	return cmd
}
