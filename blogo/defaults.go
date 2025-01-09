package blogo

import (
	"errors"
	"fmt"
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
	if d, ok := f.(fs.ReadDirFile); ok {
		return p.renderDirectory(d, w)
	}

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		_, err = w.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *defaultRenderer) renderDirectory(f fs.ReadDirFile, w io.Writer) error {
	es, err := f.ReadDir(-1)
	if err != nil {
		return err
	}

	for _, e := range es {
		_, err := w.Write([]byte(fmt.Sprintf("%s\n", e.Name())))
		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to write directory file list, file %s", e.Name()),
				err,
			)
		}
	}

	return nil
}