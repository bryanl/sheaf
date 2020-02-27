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

// NewAddManifestCommand creates an add manifest command.
func NewAddManifestCommand() *cobra.Command {
	var files []string
	var force bool

	cmd := &cobra.Command{
		Use:   "add-manifest",
		Short: "add manifest to bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires bundle path")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ma, err := sheaf.NewManifestAdder(args[0],
				sheaf.ManifestAdderForce(force))
			if err != nil {
				return err
			}
			return ma.Add(files)
		},
	}

	cmd.Flags().StringSliceVarP(&files, "filename", "f", nil,
		"filename to add (can specify multiple times)")
	cmd.Flags().BoolVar(&force, "force", false, "force (overwrite)")

	return cmd
}
