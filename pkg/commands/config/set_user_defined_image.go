/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package config

import (
	"github.com/spf13/cobra"

	"github.com/bryanl/sheaf/pkg/option"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

// NewSetUserDefinedImage creates a set user defined image command.
func NewSetUserDefinedImage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set-udi",
		Short: "Set user defined image in bundle",
		Long: `The --json-path flag supports the following BNF syntax and semantics.

Syntax

<path> ::= <identity> | <root> <subpath> | <subpath>
<identity> := ""                                         ; the current node
<root> ::= "$"                                           ; the root node of a document
<subpath> ::= <identity> | <child> <subpath> |
			  <child> <array subscript> <subpath> |
			  <recursive descent> <subpath>

<child> ::= <dot child> | <bracket child>
<dot child> ::= "." <child name> | ".*"                  ; child named <child name> or all children
<bracket child> ::= "['" <child name> "']"               ; child named <child name>

<recursive descent> ::= ".." <child name>                ; all the descendants named <child name>

<array subscript> ::= "[" <index> "]"                    ; zero or more elements of a sequence
<index> ::= <integer> | <range> | "*"                    ; specific index, range of indices, or all indices
<range> ::= <integer> ":" <integer> |                    ; start (inclusive) to end (exclusive)
            <integer> ":" <integer> ":" <integer>        ; start (inclusive) to end (exclusive) by step

Semantics

A path is logically a series of matchers that may be applied to certain manifests. To start with, the first
matcher is applied to a slice consisting of just the input document. Each matcher is applied in turn to the
slice of nodes found so far and the results are combined into a single slice, which then passes to the next
matcher, and so on. If a matcher produces an empty slice, then each subsequent matcher also produces an
empty slice and result is an empty slice.

The following matchers, with corresponding concrete syntax, are supported. See the BNF syntax above for
details of the concrete syntax.

Root: $

This matches the root node of the input YAML node. This matcher may be specified only at the start of the
path. It is optional and, if omitted, the root node is matched before the rest of the path is applied.
The output slice consists of just the root node.

Child: .childname or ['childname']

This matches the children with the given name of all the mapping nodes in the input slice. The output
slice consists of all those children. The given name may be a single child name (no periods) or a series
of single child names separated by periods. Non-mapping nodes in the input slice are not matched.

Although either form .childname or ['childname'] accepts a child name with embedded spaces, the ['childname']
form may be more convenient.

Recursive Descent: ..childname or ..*

A matcher of the form ..childname selects all the descendents of the nodes in the input slice (including
those nodes) with the given name (using the same rules as the child matcher). The output slice consists
of all the matching descendents.

A matcher of the form ..* selects all the descendents of the nodes in the input slice (including those nodes).

Array Subscript: [integer], [start:end], [start:end:step], or [*]

This matches subsequences of all the sequence nodes in the input slice. Non-sequence nodes in the input
slice are not matched.

A matcher of the form [integer] selects the corresponding node in each sequence node, with 0 meaning the first
node in the sequence, 1 the second node, and so on. A special index of -1 selects the last node in each sequence.

A matcher of the form [start:end] or [start:end:step] selects the corresponding nodes in each sequence node
starting from the start of the range (inclusive) to the end of the range (exclusive) with an optional step value
(which defaults to 1). A step value of -1 may be used to step backwards from the end of the sequence to the
start.

A matcher of the form [*] selects all the nodes in each sequence node.
`,
		Args: cobra.NoArgs,
	}

	setupSetUDI(cmd)
	return cmd
}

func setupSetUDI(cmd *cobra.Command) {
	g := option.NewGenerator(cmd, sheaf.ConfigSetUDI, "config-set-udi")
	g.WithBundlePath()
	g.WithUserDefinedImage()
}
