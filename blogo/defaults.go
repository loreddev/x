package blogo

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
)

type painTextRenderer struct{}

func NewPlainTextRenderer() Plugin {
	return &painTextRenderer{}
}

func (p *painTextRenderer) Name() string {
	return "blogo-renderer-plaintext"
}

func (p *painTextRenderer) Render(f fs.File, w io.Writer) error {
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

func (p *painTextRenderer) renderDirectory(f fs.ReadDirFile, w io.Writer) error {
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
