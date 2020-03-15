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

	"github.com/bryanl/sheaf/internal/stringutil"
	"github.com/bryanl/sheaf/pkg/sheaf"
)

func sheafInit(t *testing.T, h *harness, name, wd string) *bundle {
	err := h.runSheaf(wd, defaultSheafRunSettings, "init", name)
	require.NoError(t, err, "initialize sheaf bundle")

	return newBundle(t, filepath.Join(wd, name), h)
}

type bundle struct {
	t       *testing.T
	dir     string
	harness *harness
}

func newBundle(t *testing.T, dir string, h *harness) *bundle {
	b := bundle{
		t:       t,
		dir:     dir,
		harness: h,
	}

	return &b
}

func (b bundle) readConfig() sheaf.BundleConfig {
	var config sheaf.BundleConfig
	readJSONFile(b.t, b.configFile(), &config)
	return config
}

func (b bundle) updateConfig(fn func(config *sheaf.BundleConfig)) {
	var config sheaf.BundleConfig
	readJSONFile(b.t, b.configFile(), &config)

	fn(&config)
	writeJSONFile(b.t, b.configFile(), config)
}

func (b bundle) configFile() string {
	return filepath.Join(b.dir, "bundle.json")
}

func (b bundle) pathJoin(parts ...string) string {
	return filepath.Join(append([]string{b.dir}, parts...)...)
}

func (b bundle) archiveName() string {
	config := b.readConfig()
	return fmt.Sprintf("%s-%s.tgz", config.Name, config.Version)
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
	return fmt.Sprintf("%s%s", r.Port(t), path)
}

func (r *registry) Port(t *testing.T) string {
	require.True(t, r.started, "registry has not been started")
	cmd := exec.Command("docker", "inspect",
		"--format='{{range $p, $conf := .NetworkSettings.Ports}}{{(index $conf 0).HostIP}}:{{(index $conf 0).HostPort}}{{end}}'",
		r.id)
	data, err := cmd.CombinedOutput()
	require.NoError(t, err, "retrieve registry port: %s", string(data))

	port := string(bytes.TrimSpace(data))
	if port[0] == '\'' {
		port = port[1:]
	}
	if i := len(port) - 1; port[i] == '\'' {
		port = port[:i]
	}

	return port
}

func genRegistryPath(options wdOptions) string {
	imageName := fmt.Sprintf("%s:v1",
		stringutil.RandomWithCharset(6, stringutil.LowerAlphaCharset))

	return fmt.Sprintf("%s/%s", genRegistryRoot(options), imageName)
}

func genRegistryRoot(options wdOptions) string {
	root := fmt.Sprintf("/%s",
		stringutil.RandomWithCharset(6, stringutil.LowerAlphaCharset))

	return fmt.Sprintf("%s%s", options.registry, root)

}
