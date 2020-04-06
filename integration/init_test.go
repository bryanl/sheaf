/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_sheaf_init(t *testing.T) {
	bundleName := "my-bundle"

	withWorkingDirectory(t, func(options wdOptions) {

		cases := []struct {
			name          string
			bundleName    string
			args          []string
			runnerOptions []sheafInitRunnerOption
		}{
			{
				name:       "with no options",
				bundleName: bundleName,
			},
			{
				name:       "with a bundle directory",
				bundleName: bundleName,
				args:       []string{"--bundle-path", "custom-dir"},
				runnerOptions: []sheafInitRunnerOption{
					func(initRunner sheafInitRunner) sheafInitRunner {
						initRunner.bundlePath = "custom-dir"
						return initRunner
					},
				},
			},
		}

		for _, tc := range cases {
			sir := newSheafInitRunner(t, testHarness, tc.name, options.dir, tc.bundleName, tc.runnerOptions...)
			sir.Run(tc.args...)
		}
	})
}

type sheafInitRunnerOption func(initRunner sheafInitRunner) sheafInitRunner

type sheafInitRunner struct {
	name             string
	workingDirectory string
	bundleName       string
	bundlePath       string
	harness          *harness
	t                *testing.T
}

func newSheafInitRunner(t *testing.T, h *harness, name, wd, bundleName string, options ...sheafInitRunnerOption) *sheafInitRunner {
	sir := sheafInitRunner{
		t:                t,
		name:             name,
		workingDirectory: wd,
		bundleName:       bundleName,
		harness:          h,
	}

	require.NotEmpty(t, name)
	require.NotEmpty(t, wd)
	require.NotEmpty(t, bundleName)

	for _, option := range options {
		sir = option(sir)
	}

	return &sir
}

func (r *sheafInitRunner) Run(args ...string) {
	r.t.Run(r.name, func(t *testing.T) {
		args = append([]string{"init", "--bundle-name", r.bundleName}, args...)

		root := filepath.Join(r.workingDirectory, r.bundleName)
		if r.bundlePath != "" {
			root = filepath.Join(r.workingDirectory, r.bundlePath)
		}

		_, err := r.harness.runSheaf(r.workingDirectory, args...)
		require.NoError(t, err)

		t.Run("creates a bundle directory", func(t *testing.T) {
			checkDirExists(t, root)
		})

		t.Run("creates a manifests directory", func(t *testing.T) {
			checkDirExists(t, filepath.Join(root, "app", "manifests"))
		})

		t.Run("creates a bundle configuration", func(t *testing.T) {
			checkFileExists(t, filepath.Join(root, "bundle.json"))
		})
	})
}
