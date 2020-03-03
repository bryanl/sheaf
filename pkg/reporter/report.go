/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package reporter

var (
	// Default is the default reporter.
	Default = newStd()
)

// Reporter reports information to sheaf users.
type Reporter interface {
	// Header prints a header.
	Header(string)
	// Report creates report output.
	Report(format string)
	// Reportf creates formatted report output.
	Reportf(format string, a ...interface{})
}
