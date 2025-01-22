// Copyright 2025-present Gustavo "Guz" L. de Mello
// Copyright 2025-present The Lored.dev Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fs

import (
	iofs "io/fs"
)

// Wraps the provided [iofs.FS] file system so it can be used as a file system for blogo.
// [Metadata] from this [FS] will be empty, and by default, mutable.
func FromIOFS(fsys iofs.FS, immutable ...bool) FS {
	if fsys == nil {
		return nil
	}

	m := MetadataMap(map[string]any{})
	i := false
	if len(immutable) > 0 && immutable[0] {
		m = ImmutableMetadata(m)
		i = true
	}

	return &wrapperFS{
		fsys:      fsys,
		metadata:  m,
		immutable: i,
	}
}

type wrapperFS struct {
	fsys      iofs.FS
	metadata  Metadata
	immutable bool
}

func (f *wrapperFS) Metadata() Metadata {
	return f.metadata
}

func (f *wrapperFS) Open(name string) (File, error) {
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	return FromIOFile(file, f.immutable), nil
}

// Wraps the provided [iofs.File] so it can be used as a file system for blogo.
// [Metadata] from this [File] will be empty, and by default, mutable.
func FromIOFile(file iofs.File, immutable ...bool) File {
	m := MetadataMap(map[string]any{})
	if len(immutable) > 0 && immutable[0] {
		m = ImmutableMetadata(m)
	}

	return &wrapperFile{
		File:     file,
		metadata: m,
	}
}

type wrapperFile struct {
	iofs.File
	metadata Metadata
}

func (f *wrapperFile) Metadata() Metadata {
	return f.metadata
}
