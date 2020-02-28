/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package codec

import "github.com/bryanl/sheaf/pkg/sheaf"

var (
	// Default is the default JSON codec.
	Default = &JSONCodec{
		JSONEncoder: DefaultEncoder,
		JSONDecoder: DefaultDecoder,
	}
)

// JSONCodec is a JSON codec.
type JSONCodec struct {
	*JSONEncoder
	*JSONDecoder
}

var _ sheaf.Codec = &JSONCodec{}
