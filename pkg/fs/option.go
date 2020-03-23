/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"github.com/bryanl/sheaf/pkg/sheaf"
)

type options struct {
	bundleConfigCodec sheaf.BundleConfigCodec
}

// Option is a functional option for configuring fs
type Option func(options *options)

func makeCreateBundleOptions(optionList ...Option) options {
	opts := options{
		bundleConfigCodec: &BundleConfigCodec{},
	}

	for _, o := range optionList {
		o(&opts)
	}

	return opts
}
