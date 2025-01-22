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

type FS interface {
	Metadata() Metadata
	Open(name string) (File, error)
}

type File interface {
	iofs.File
	Metadata() Metadata
}

// Alias for "io/fs"
var (
	ErrInvalid    = iofs.ErrInvalid    // "invalid argument"
	ErrPermission = iofs.ErrPermission // "permission denied"
	ErrExist      = iofs.ErrExist      // "file already exists"
	ErrNotExist   = iofs.ErrNotExist   // "file does not exist"
	ErrClosed     = iofs.ErrClosed     // "file already closed"
)

// Alias for "io/fs"
func FormatDirEntry(dir DirEntry) string { return iofs.FormatDirEntry(dir) }

// Alias for "io/fs"
func FormatFileInfo(info FileInfo) string { return iofs.FormatFileInfo(info) }

// TODO: func Glob(fsys FS, pattern string) (matches []string, err error) { return iofs.Glob(fsys, pattern) }
// TODO: func ReadFile(fsys FS, name string) ([]byte, error) {return iofs.ReadFile(fsys, name)}

// Alias for "io/fs"
func ValidPath(name string) bool { return iofs.ValidPath(name) }

// TODO: func WalkDir(fsys FS, root string, fn WalkDirFunc) error { return iofs.WalkDir(fsys, root, fn) }

// Alias for "io/fs"
type (
	DirEntry    = iofs.DirEntry
	FileInfo    = iofs.FileInfo
	FileMode    = iofs.FileMode
	GlobFS      = iofs.GlobFS
	PathError   = iofs.PathError
	ReadDirFS   = iofs.ReadDirFS
	ReadDirFile = iofs.ReadDirFile
	ReadFileFS  = iofs.ReadFileFS
	StatFS      = iofs.StatFS
	SubFS       = iofs.StatFS
	WalkDirFunc = iofs.WalkDirFunc
)
