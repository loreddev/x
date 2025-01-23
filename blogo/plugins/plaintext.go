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

package plugins

import (
	"errors"
	"fmt"
	"io"
	"io/fs"

	"forge.capytal.company/loreddev/x/blogo/plugin"
)

const plainTextName = "blogo-plaintext-renderer"

type painText struct{}

func NewPlainText() plugin.Renderer {
	return &painText{}
}

func (p *painText) Name() string {
	return plainTextName
}

func (p *painText) Render(f fs.File, w io.Writer) error {
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

func (p *painText) renderDirectory(f fs.ReadDirFile, w io.Writer) error {
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
