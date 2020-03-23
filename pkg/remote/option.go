/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package remote

type options struct {
	insecureRegistry bool
}

// Option is a functional option for configuring remote.
type Option func(o *options)

// WithInsecureRegistry sets insecure registry.
func WithInsecureRegistry(forceInsecure bool) Option {
	return func(o *options) {
		o.insecureRegistry = forceInsecure
	}
}
