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
	"fmt"
	"log"
	"os"
)

// Unpacker unpacks bundles.
type Unpacker struct {
	ArchivePath string
	Dest        string
}

// UnpackerOption is an option for configuration Unpacker.
type UnpackerOption func(u Unpacker) Unpacker

// UnpackerArchivePath sets the archive path for Unpacker.
func UnpackerArchivePath(p string) UnpackerOption {
	return func(u Unpacker) Unpacker {
		u.ArchivePath = p
		return u
	}
}

// UnpackerDest sets the destination for Unpacker.
func UnpackerDest(p string) UnpackerOption {
	return func(u Unpacker) Unpacker {
		u.Dest = p
		return u
	}
}

// NewUnpacker creates an instance of Unpacker.
func NewUnpacker(options ...UnpackerOption) *Unpacker {
	u := Unpacker{}

	for _, option := range options {
		u = option(u)
	}

	return &u
}

// Unpack unpacks an archive.
func (u *Unpacker) Unpack() error {
	if u.ArchivePath == "" {
		return fmt.Errorf("archive path is blank")
	}

	if u.Dest == "" {
		return fmt.Errorf("destination is blank")
	}

	f, err := os.Open(u.ArchivePath)
	if err != nil {
		return fmt.Errorf("open %q: %w", u.ArchivePath, err)
	}

	defer func() {
		if cErr := f.Close(); cErr != nil {
			log.Printf("close %q: %v", u.ArchivePath, err)
		}
	}()

	if err := Unarchive(f, u.Dest); err != nil {
		return fmt.Errorf("unarchive %q: %w", u.ArchivePath, err)
	}

	return nil
}
