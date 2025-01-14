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

package blogo

import (
	"io/fs"
)

type FS interface {
	Metadata() Metadata
	Open(name string) (File, error)
}

type File interface {
	fs.File
	Metadata() Metadata
}

type wrapperFS struct {
	fs.FS
	metadata  Metadata
	immutable bool
}

func FsFS(f fs.FS, immutable ...bool) FS {
	var m Metadata
	var i bool
	if len(immutable) > 0 && immutable[0] {
		i = true
		m = ImmutableMetadata(MetadataMap(map[string]any{}))
	} else {
		i = false
		m = MetadataMap(map[string]any{})
	}

	return &wrapperFS{
		FS:        f,
		metadata:  m,
		immutable: i,
	}
}

func (f *wrapperFS) Metadata() Metadata {
	return f.metadata
}

func (f *wrapperFS) Open(name string) (File, error) {
	file, err := f.FS.Open(name)
	if err != nil {
		return nil, err
	}
	return FsFile(file, f.immutable), nil
}

type wrapperFile struct {
	fs.File
	metadata Metadata
}

func FsFile(f fs.File, immutable ...bool) File {
	var m Metadata
	if len(immutable) > 0 && immutable[0] {
		m = ImmutableMetadata(MetadataMap(map[string]any{}))
	} else {
		m = MetadataMap(map[string]any{})
	}

	return &wrapperFile{
		File:     f,
		metadata: m,
	}
}

func (f *wrapperFile) Metadata() Metadata {
	return f.metadata
}
