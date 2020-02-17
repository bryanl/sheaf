/*
 * Copyright 2020 Sheaf Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sheaf

import (
	"sort"

	"github.com/pivotal/image-relocation/pkg/image"
)

// imageSlice attaches the methods of sort.Interface to []image.Name, sorting in increasing string order.
type imageSlice []image.Name

func (p imageSlice) Len() int           { return len(p) }
func (p imageSlice) Less(i, j int) bool { return p[i].String() < p[j].String() }
func (p imageSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p imageSlice) Sort() { sort.Sort(p) }

func sortImages(a []image.Name) { imageSlice(a).Sort() }
