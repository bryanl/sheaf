/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/fs"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewSetUserDefinedImage creates a set user defined image command.
func NewSetUserDefinedImage() *cobra.Command {
	var bundlePath string
	var apiVersion string
	var kind string
	var jsonPath string
	var imageType string

	cmd := &cobra.Command{
		Use:   "set-udi",
		Short: "Set user defined image in bundle",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			a := fs.NewUserDefinedImageSetter()

			udi := sheaf.UserDefinedImage{
				APIVersion: apiVersion,
				Kind:       kind,
				JSONPath:   jsonPath,
				Type:       sheaf.UserDefinedImageType(imageType),
			}

			return a.Set(bundlePath, udi)
		},
	}

	cmd.Flags().StringVar(&apiVersion, "api-version", "", "api version")
	cmd.Flags().StringVar(&kind, "kind", "", "kind")
	cmd.Flags().StringVar(&jsonPath, "json-path", "", "kind")
	cmd.Flags().StringVar(&imageType, "type", string(sheaf.SingleResult),
		fmt.Sprintf("type of user defined image (valid types: %s)",
			strings.Join(sheaf.UserDefinedImageTypes, ",")))

	cmd.Flags().StringVar(&bundlePath, "bundle-path", ".", "bundle path")

	return cmd
}
