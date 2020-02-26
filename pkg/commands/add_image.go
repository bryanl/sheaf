/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
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

// NewAddImageCommand creates an add image command.
func NewAddImageCommand() *cobra.Command {
	var images []string

	cmd := &cobra.Command{
		Use:   "add-image",
		Short: "add image to bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires bundle path")
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
