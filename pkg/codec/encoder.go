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
	// DefaultEncoder is the default encoder.
	DefaultEncoder = &JSONEncoder{}
)

// JSONEncoder encodes a value to JSON.
type JSONEncoder struct{}

var _ sheaf.Encoder = &JSONEncoder{}

// Encode encodes a value to JSON.
func (e JSONEncoder) Encode(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}
