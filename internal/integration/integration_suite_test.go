// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"go.uber.org/multierr"
)

var (
	testHarness *harness
)

func TestMain(m *testing.M) {
	code, err := run(m)
	if err != nil {
		log.Printf("%v", err)
		os.Exit(code)
	}

	os.Exit(code)
}

func run(m *testing.M) (code int, err error) {
	dir, err := ioutil.TempDir("", "build-dir")
	if err != nil {
		return 1, fmt.Errorf("create build directory: %w", err)
	}

	defer func() {
		errors := []error{err}
		if cErr := testHarness.cleanup(); cErr != nil {
			errors = append(errors, fmt.Errorf("cleanup: %w", cErr))
		}

		err = multierr.Combine(errors...)
	}()

	testHarness, err = newRunner(dir)
	if err != nil {
		return 1, fmt.Errorf("creat test harness: %w", err)
	}

	code = m.Run()
	if code != 0 {
		return code, fmt.Errorf("non zero error code %d", code)
	}

	return code, nil
}
