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

// NewShowManifests creates a "show manifests" command.
func NewShowManifests() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-manifests",
		Short: "Prints manifest to standard out",
		Args:  cobra.NoArgs,
	}

	setupShowManifest(cmd)
	return cmd
}

func setupShowManifest(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ArchiveShowManifests, "archive-show-manifests")
	g.WithBundlePath()
	g.WithArchive()
	g.WithPrefix()
}
