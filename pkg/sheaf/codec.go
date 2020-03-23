/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

//go:generate mockgen -destination=../mocks/mock_decoder.go -package mocks github.com/bryanl/sheaf/pkg/sheaf Decoder
//go:generate mockgen -destination=../mocks/mock_encoder.go -package mocks github.com/bryanl/sheaf/pkg/sheaf Encoder
//go:generate mockgen -destination=../mocks/mock_codec.go -package mocks github.com/bryanl/sheaf/pkg/sheaf Codec

// Decoder decodes bytes into a value.
type Decoder interface {
	Decode([]byte, interface{}) error
}

// Encoder encodes a value into bytes.
type Encoder interface {
	Encode(interface{}) ([]byte, error)
}

// Codec combines Decoder and Encoder
type Codec interface {
	Decoder
	Encoder
}
