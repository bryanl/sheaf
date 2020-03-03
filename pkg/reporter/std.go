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

type stdOption func(std) std

func stdOptionOutput(w io.Writer) stdOption {
	return func(std std) std {
		std.w = w
		return std
	}
}

type std struct {
	w      io.Writer
	indent string
}

var _ Reporter = &std{}

func newStd(options ...stdOption) *std {
	r := std{
		w:      os.Stdout,
		indent: defaultIndent,
	}

	for _, option := range options {
		r = option(r)
	}

	return &r
}

func (r std) Header(text string) {
	fmt.Fprintln(r.w, text)
}

func (r std) Report(text string) {
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

func (r std) Reportf(format string, a ...interface{}) {
	r.Report(fmt.Sprintf(format, a...))
}
