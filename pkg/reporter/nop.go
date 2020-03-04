/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package reporter

// Nop is a no-op reporter.
type Nop struct{}

var _ Reporter = &Nop{}

// Header implements Reporter#Header.
func (n Nop) Header(string) {
}

// Headerf implements Reporter#Headerf.
func (n Nop) Headerf(format string, a ...interface{}) {
}

// Report implements Reporter#Report.
func (n Nop) Report(format string) {
}

// Reportf implements Reporter#Reportf.
func (n Nop) Reportf(format string, a ...interface{}) {
}
