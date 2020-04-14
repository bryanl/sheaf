/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type harness struct {
	buildDir string
	sheafBin string
}

func newRunner(buildDir string) (*harness, error) {
	r := harness{
		buildDir: buildDir,
	}

	if err := r.buildSheaf(); err != nil {
		return nil, fmt.Errorf("build sheaf: %w", err)
	}

	return &r, nil
}

func (r *harness) buildSheaf() error {
	if runtime.GOOS == "windows" {
		r.sheafBin = filepath.Join(r.buildDir, "sheaf.exe")
	} else {
		r.sheafBin = filepath.Join(r.buildDir, "sheaf")
	}

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

	return nil
}

type sheafOutput struct {
	Stdout bytes.Buffer
	Stderr bytes.Buffer
}

func (r harness) runSheaf(workingDirectory string, args ...string) (*sheafOutput, error) {
	cmd := exec.Command(r.sheafBin, args...)

	var stdout, stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = workingDirectory

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run sheaf: %w", err)
	}

	return &sheafOutput{
		Stdout: stdout,
		Stderr: stderr,
	}, nil
}

func (r harness) cleanup() error {
	if r.buildDir != "" {
		if err := os.RemoveAll(r.buildDir); err != nil {
			return err
		}
	}

	return nil
}
