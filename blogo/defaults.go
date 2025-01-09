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

type defaultRenderer struct{}

func (p *defaultRenderer) Name() string {
	return "blogo-default-renderer"
}

func (p *defaultRenderer) Render(f fs.File, w io.Writer) error {
	return nil
}
