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

// Alias for "io/fs"
var (
	ErrInvalid    = iofs.ErrInvalid    // "invalid argument"
	ErrPermission = iofs.ErrPermission // "permission denied"
	ErrExist      = iofs.ErrExist      // "file already exists"
	ErrNotExist   = iofs.ErrNotExist   // "file does not exist"
	ErrClosed     = iofs.ErrClosed     // "file already closed"
)

// Provides access to a hierarchical file system, similar to [iofs.FS]. Implementations
// can load files, fetch them from network, read from disk, as they are opened/on demand
// if they seem fit to do so.
//
// A file system may implement additional interfaces such as [ReadFileFS].
type FS interface {
	// Returns [Metadata] about the file system.
	//
	// Implementations should return a empty [Metadata] instead of a nil value for when
	// the file system doesn't have any additional metadata about it or when a error
	// occurs when getting said metadata.
	//
	// [plugin.Sourcer] may add prefixes to their metadata keys.
	Metadata() Metadata
	// Open, similar to [iofs.File.Open], opens the named file.
	//
	// When it returns an error, it should be of type [*PathError] with the Op field
	// set to "open", the Path field set to name, and the Err field describing the problem.
	//
	// It should reject attempts of opening names that do not satisfy [ValidPath], returning
	// a [*PathError] with Err set to [ErrInvalid] or [ErrNotExist]
	//
	// Implementations may find the file on demand or fetch them via http request, depending on
	// the underlying source of the file.
	Open(name string) (File, error)
}

// Provides access to a single file, similar to [iofs.FS]. Implementations may read the file
// on demand and/or fetch it's contents and data from a outside source such as via HTTP requests,
// depending on the underlying source of the file.
//
// Directory files should also implement [ReadDirFile]. A file may implement [io.ReaderAt] or
// [io.Seeker] as optimizations.
type File interface {
	// Returns [Metadata] about the file.
	//
	// Implementations should return a empty [Metadata] instead of a nil value for when
	// the file doesn't have any additional metadata about it or when a error occurs while getting
	// said metadata.
	//
	// [plugin.Sourcer] may add prefixes to their metadata keys.
	Metadata() Metadata

	Stat() (FileInfo, error)
	Read([]byte) (int, error)
	Close() error
}

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
