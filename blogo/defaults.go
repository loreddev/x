package blogo

import (
	"io"
	"io/fs"
)

type defaultSourcer struct{}

func (p *defaultSourcer) Name() string {
	return "blogo-defaults-sourcer"
}

func (p *defaultSourcer) Source() (fs.FS, error) {
	return emptyFS{}, nil
}

type emptyFS struct{}

func (f emptyFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

