/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import "github.com/bryanl/sheaf/pkg/sheaf"

// UserDefinedImageDeleterOption is a functional option for UserDefinedImageDeleter.
type UserDefinedImageDeleterOption func(d UserDefinedImageDeleter) UserDefinedImageDeleter

// UserDefinedImageDeleter deletes user defined images.
type UserDefinedImageDeleter struct {
	bundleFactory      sheaf.BundleFactory
	bundleConfigWriter sheaf.BundleConfigWriter
}

// NewUserDefinedImageDeleter creates an instance of UserDefinedImageDeleter.
func NewUserDefinedImageDeleter(options ...UserDefinedImageDeleterOption) *UserDefinedImageDeleter {
	d := UserDefinedImageDeleter{
		bundleFactory:      DefaultBundleFactory,
		bundleConfigWriter: DefaultBundleConfigWriter,
	}

	for _, option := range options {
		d = option(d)
	}

	return &d
}

// Delete deletes an image using a key.
func (d UserDefinedImageDeleter) Delete(uri string, key sheaf.UserDefinedImageKey) error {
	b, err := d.bundleFactory(uri)
	if err != nil {
		return err
	}

	config := updateUDI(b.Config(), func(u udiMap) {
		delete(u, key)
	})

	return d.bundleConfigWriter(b, config)
}
