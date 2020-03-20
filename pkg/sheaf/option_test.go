/*
 * Copyright 2020 Sheaf Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package sheaf_test

import (
	"fmt"

	"github.com/golang/mock/gomock"

	"github.com/bryanl/sheaf/pkg/sheaf"
)

type bundleFactoryFunc func(controller *gomock.Controller) sheaf.BundleFactoryFunc

func genManifest(name string) sheaf.BundleManifest {
	return sheaf.BundleManifest{
		ID:   name,
		Data: []byte(fmt.Sprintf("file: %s", name)),
	}
}
