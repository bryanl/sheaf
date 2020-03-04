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
	// Header prints a header from a formatted string.
	Headerf(format string, a ...interface{})
	// Report creates report output.
	Report(format string)
	// Reportf creates formatted report output.
	Reportf(format string, a ...interface{})
}
