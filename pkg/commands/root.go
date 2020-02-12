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

import "github.com/spf13/cobra"

// NewRootCommand creates a root command.
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "sheaf",
		Short:        "sheaf bundles Kubernetes applications",
		SilenceUsage: true,
	}

	cmd.AddCommand(
		NewInitCommand(),
		NewAddManifestCommand(),
		NewAddImageCommand(),
		NewPackCommand(),
		NewStageCommand(),
		NewGenManifestCommand())

	return cmd
}

// Execute executes the root command for sheaf.
func Execute() error {
	root := NewRootCommand()
	if err := root.Execute(); err != nil {
		return err
	}

	return nil
}
