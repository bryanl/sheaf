/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package stringutil

import (
	"math/rand"
	"time"
)

const (
	// LowerAlphaCharset is a charset of lower case letters.
	LowerAlphaCharset = "abcedefghijklmnopqrstuvwxyz"
	// UpperAlphaCharset is a charset of upper case letters.
	UpperAlphaCharset = "ABCDEFGHIGJKLMNOPQRSTUVWXYZ"
	// NumberCharset is a charset of numbers.
	NumberCharset = "0123456789"
	// DefaultChartSet is charset of lower case letters, upper case letters,
	// and numbers.
	DefaultChartSet = LowerAlphaCharset + UpperAlphaCharset + NumberCharset
)

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// RandomWithCharset generates a random string with length given a charset.
func RandomWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// Random generates a random string with length using the default charset.
func Random(length int) string {
	return RandomWithCharset(length, DefaultChartSet)
}
