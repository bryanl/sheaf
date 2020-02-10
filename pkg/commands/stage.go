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

// NewStageCommand creates a stage command.
func NewStageCommand() *cobra.Command {
	var unpackDir string

	cmd := &cobra.Command{
		Use:   "stage",
		Short: "stage a bundle",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("requires bundle location and registry prefix")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			config := sheaf.StageConfig{
				ArchivePath:    args[0],
				RegistryPrefix: args[1],
				UnpackDir:      unpackDir,
			}

			return sheaf.Stage(config)
		},
	}

	cmd.Flags().StringVar(&unpackDir, "unpack-dir", "", "directory to unpack bundle to")

	return cmd
}
