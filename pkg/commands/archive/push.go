/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package archive

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/option"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewPushCommand create a push command.
func NewPushCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "push",
		// TODO: support relocating images
		Short: "Push sheaf archive to registry",
		Args:  cobra.NoArgs,
	}

	setupPush(cmd)
	return cmd
}

func setupPush(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ArchivePush, "archive-push")
	g.WithBundlePath()
	g.WithArchive()
	g.WithReference()
	g.WithInsecureRegistry()

}
