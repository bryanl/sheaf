// +build integration

/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package integration_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

func sheafInit(t *testing.T, h *harness, name, wd string) *bundle {
	err := h.runSheaf(wd, defaultSheafRunSettings, "init", name)
	require.NoError(t, err, "initialize sheaf bundle")

	b := bundle{
		dir:     filepath.Join(wd, name),
		harness: h,
	}

	return &b
}

type bundle struct {
	dir     string
	harness *harness
}

func (b bundle) readConfig(t *testing.T) sheaf.BundleConfig {
	var config sheaf.BundleConfig
	readJSONFile(t, b.configFile(), &config)
	return config
}

func (b bundle) updateConfig(t *testing.T, fn func(config *sheaf.BundleConfig)) {
	var config sheaf.BundleConfig
	readJSONFile(t, b.configFile(), &config)

	fn(&config)
	writeJSONFile(t, b.configFile(), config)
}

func (b bundle) configFile() string {
	return filepath.Join(b.dir, "bundle.json")
}

func (b bundle) pathJoin(parts ...string) string {
	return filepath.Join(append([]string{b.dir}, parts...)...)
}

type registry struct {
	id      string
	started bool

	mu sync.Mutex
}

func newRegistry() *registry {
	r := registry{}
	return &r
}

func (r *registry) Start(t *testing.T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	require.False(t, r.started, "registry has already been started")

	cmd := exec.Command("docker", "run", "-p", "5000", "--rm", "-d", "registry:2")
	data, err := cmd.Output()
	require.NoError(t, err, "start registry: %s", string(data))

	r.id = string(bytes.TrimSpace(data))
	r.started = true
}

func (r *registry) Stop(t *testing.T) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.started {
		return
	}

	cmd := exec.Command("docker", "stop", r.id)
	require.NoError(t, cmd.Run(), "stop registry")
	r.started = false
}

func (r *registry) Ref(t *testing.T, path string) string {
	require.True(t, r.started, "registry has not been started")
	return fmt.Sprintf("localhost:%s%s", r.port(t), path)
}

func (r *registry) port(t *testing.T) string {
	cmd := exec.Command("docker", "inspect",
		"--format='{{range $p, $conf := .NetworkSettings.Ports}}{{(index $conf 0).HostPort}}{{end}}'",
		r.id)
	data, err := cmd.Output()
	require.NoError(t, err, "retrieve registry port")

	port := string(bytes.TrimSpace(data))
	if port[0] == '\'' {
		port = port[1:]
	}
	if i := len(port) - 1; port[i] == '\'' {
		port = port[:i]
	}

	return port
}
