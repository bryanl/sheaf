/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package reporter

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	defaultIndent = "  "
)

// StdOption is a functional option for configuring Std.
type StdOption func(Std) Std

// WithWriter sets the writer.
func WithWriter(w io.Writer) StdOption {
	return func(std Std) Std {
		std.w = w
		return std
	}
}

// Std is the standard reporter.
type Std struct {
	w      io.Writer
	indent string
}

var _ Reporter = &Std{}

// New creates an instance of Std.
func New(options ...StdOption) *Std {
	r := Std{
		w:      os.Stdout,
		indent: defaultIndent,
	}

	for _, option := range options {
		r = option(r)
	}

	return &r
}

// Header prints a header.
func (r Std) Header(text string) {
	fmt.Fprintln(r.w, text)
}

// Headerf prints a formatted header.
func (r Std) Headerf(format string, a ...interface{}) {
	r.Header(fmt.Sprintf(format, a...))
}

// Report prints a report. Reports are indented.
func (r Std) Report(text string) {
	if text[len(text)-1:] == "\n" {
		result := ""
		for _, s := range strings.Split(text[:len(text)-1], "\n") {
			result += r.indent + s + "\n"
		}
	}
	result := ""
	for _, s := range strings.Split(strings.TrimRight(text, "\n"), "\n") {
		result += r.indent + s + "\n"
	}

	fmt.Fprint(r.w, result[:len(result)-1])
	fmt.Fprintln(r.w)
}

// Reportf prints a formatted report. Reports are indented.
func (r Std) Reportf(format string, a ...interface{}) {
	r.Report(fmt.Sprintf(format, a...))
}
