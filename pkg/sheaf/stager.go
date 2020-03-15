/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

// Stager is the interface that wraps the Stage function.
type Stager interface {
	// Stage stages an archive by URI to a registry given a prefix.
	Stage(archiveURI, registryPrefix string, insecureRegistry bool) error
}
