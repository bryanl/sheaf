/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package fs

import (
	"fmt"
	"sort"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

// UserDefinedImageSetterOption is a functional option for configuring UserDefinedImageSetter.
type UserDefinedImageSetterOption func(s UserDefinedImageSetter) UserDefinedImageSetter

// UserDefinedImageSetter sets user defined images.
type UserDefinedImageSetter struct {
	bundleFactory      sheaf.BundleFactory
	BundleConfigWriter sheaf.BundleConfigWriter
}

// NewUserDefinedImageSetter creates an instance of UserDefinedImageSetter.
func NewUserDefinedImageSetter(options ...UserDefinedImageSetterOption) *UserDefinedImageSetter {
	a := UserDefinedImageSetter{
		bundleFactory:      DefaultBundleFactory,
		BundleConfigWriter: DefaultBundleConfigWriter,
	}

	for _, option := range options {
		a = option(a)
	}

	return &a
}

// Set sets user defined images in the bundle config. Images are keyed of api version and kind and duplicates
// will be overwritten.
func (s UserDefinedImageSetter) Set(uri string, udi sheaf.UserDefinedImage) error {
	if err := udi.Validate(); err != nil {
		return fmt.Errorf("validate user defined image: %w", err)
	}

	b, err := s.bundleFactory(uri)
	if err != nil {
		return err
	}

	config := updateUDI(b.Config(), func(u udiMap) {
		key := sheaf.UserDefinedImageKey{
			APIVersion: udi.APIVersion,
			Kind:       udi.Kind,
		}
		u[key] = udi
	})

	return s.BundleConfigWriter(b, config)
}

type udiMap map[sheaf.UserDefinedImageKey]sheaf.UserDefinedImage

func updateUDI(config sheaf.BundleConfig, fn func(udiMap)) sheaf.BundleConfig {
	m := udiMap{}
	for _, cur := range config.UserDefinedImages {
		key := sheaf.UserDefinedImageKey{
			APIVersion: cur.APIVersion,
			Kind:       cur.Kind,
		}

		m[key] = cur
	}

	fn(m)

	var list []sheaf.UserDefinedImage

	var keys []sheaf.UserDefinedImageKey
	for k := range m {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].APIVersion < keys[j].APIVersion {
			return true
		}
		if keys[i].APIVersion > keys[j].APIVersion {
			return false
		}
		return keys[i].Kind < keys[j].Kind
	})

	for _, key := range keys {
		list = append(list, m[key])
	}

	config.UserDefinedImages = list
	return config
}
