/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package codec

import (
	"encoding/json"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

var (
	// DefaultDecoder is the default decoder.
	DefaultDecoder = &JSONDecoder{}
)

// JSONDecoder decodes JSON.
type JSONDecoder struct {
}

var _ sheaf.Decoder = &JSONDecoder{}

// Decode decodes JSON.
func (d JSONDecoder) Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
