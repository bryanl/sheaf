/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf

import (
	"fmt"
	"os"
)

func ensureBundlePath(bundlePath string) error {
	fi, err := os.Stat(bundlePath)
	if err == nil {
		if !fi.IsDir() {
			return fmt.Errorf("is not a directory")
		}

		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("is invalid: %w", err)
	}

	return nil
}
