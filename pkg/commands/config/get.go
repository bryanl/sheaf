/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/fs"
)

// NewGetCommand creates `config get` command.
func NewGetCommand() *cobra.Command {
	cmd := cobra.Command{
		Use:   "get",
		Short: "Get bundle configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := fs.DefaultBundleFactory(".")
			if err != nil {
				return err
			}

			data, err := json.MarshalIndent(b.Config(), "", "  ")
			if err != nil {
				return err
			}

			fmt.Print(string(data))
			return nil
		},
	}

	return &cmd
}
