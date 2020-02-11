/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
