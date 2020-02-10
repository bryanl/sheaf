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
	"os"

	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewGenManifestCommand creates a gen manifest command.
func NewGenManifestCommand() *cobra.Command {
	var prefix string

	cmd := &cobra.Command{
		Use:   "gen-manifest",
		Short: "generate manifest",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("requires path to archive")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mg := sheaf.NewManifestGenerator(
				sheaf.ManifestGeneratorArchivePath(args[0]),
				sheaf.ManifestGeneratorPrefix(prefix))
			return mg.Generate(os.Stdout)
		},
	}

	cmd.Flags().StringVar(&prefix, "prefix", "", "registry prefix")

	return cmd
}
