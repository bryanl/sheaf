/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package goutil

import (
	"io"
	"log"
)

// Close closes a closer and logs a message if there is an error.
func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Printf("close: %v", err)
	}
}
