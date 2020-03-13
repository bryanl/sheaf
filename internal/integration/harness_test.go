// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type harness struct {
	buildDir string
	sheafBin string
}

func newRunner(buildDir string) (*harness, error) {
	r := harness{
		buildDir: buildDir,
		sheafBin: filepath.Join(buildDir, "sheaf"),
	}

	if err := r.buildSheaf(); err != nil {
		return nil, fmt.Errorf("build sheaf: %w", err)
	}

	return &r, nil
}

func (r *harness) buildSheaf() error {
	args := []string{
		"build",
		"-o", r.sheafBin,
		"github.com/bryanl/sheaf/cmd/sheaf",
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return err
	}

	r.sheafBin = filepath.Join(r.buildDir, "sheaf")
	return nil
}

type sheafRunSettings struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func genSheafRunSettings() sheafRunSettings {
	return sheafRunSettings{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

var (
	defaultSheafRunSettings = sheafRunSettings{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
)

func (r harness) runSheaf(workingDirectory string, settings sheafRunSettings, args ...string) error {
	cmd := exec.Command(r.sheafBin, args...)
	cmd.Stdin = settings.Stdin
	cmd.Stdout = settings.Stdout
	cmd.Stderr = settings.Stderr
	cmd.Dir = workingDirectory

	return cmd.Run()
}

func (r harness) cleanup() error {
	if r.buildDir != "" {
		if err := os.RemoveAll(r.buildDir); err != nil {
			return err
		}
	}

	return nil
}
