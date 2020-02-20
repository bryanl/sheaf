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

package images_test

import (
	"testing"

	"github.com/bryanl/sheaf/pkg/images"
	"github.com/stretchr/testify/require"
)

func TestMarshalling(t *testing.T) {
	s, err := images.New([]string{"a", "b"})
	require.NoError(t, err)

	sb, err := s.MarshalJSON()
	require.NoError(t, err)

	var u images.Set
	err = (&u).UnmarshalJSON(sb)
	require.NoError(t, err)

	require.Equal(t, s, u)
}

func TestNullUnmarshalling(t *testing.T) {
	var u images.Set
	err := (&u).UnmarshalJSON([]byte("null"))
	require.NoError(t, err)

	require.Equal(t, images.Empty, u)
}
