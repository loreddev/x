package blogo

import "io/fs"

type Plugin interface {
	Name() string
}

type SourcerPlugin interface {
	Plugin
	Source() (fs.FS, error)
}
