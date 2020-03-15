/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package codec

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

var (
	// DefaultEncoder is the default encoder.
	DefaultEncoder = &JSONEncoder{}
)

// JSONEncoder encodes a value to JSON.
type JSONEncoder struct{}

var _ sheaf.Encoder = &JSONEncoder{}

// Encode encodes a value to indented. JSON.
func (e JSONEncoder) Encode(v interface{}) ([]byte, error) {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		return nil, fmt.Errorf("json encode: %w", err)
	}

	return b.Bytes(), nil
}
