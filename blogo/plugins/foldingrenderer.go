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
	"bytes"
	"fmt"
	"io"
	"log/slog"

	"forge.capytal.company/loreddev/x/blogo/fs"
	"forge.capytal.company/loreddev/x/blogo/plugin"
)

const foldingRendererPluginName = "blogo-foldingrenderer-renderer"

type foldingRenderer struct {
	plugins []plugin.Renderer

	panicOnInit bool

	log *slog.Logger
}

type FoldingRendererOpts struct {
	PanicOnInit bool

	Logger *slog.Logger
}

type FoldingRenderer interface {
	plugin.WithPlugins
	plugin.Renderer
}

func NewFoldingRenderer(opts ...FoldingRendererOpts) FoldingRenderer {
	opt := FoldingRendererOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	opt.Logger = opt.Logger.WithGroup(foldingRendererPluginName)

	return &foldingRenderer{
		plugins: []plugin.Renderer{},

		log: opt.Logger,
	}
}

func (r *foldingRenderer) Name() string {
	return foldingRendererPluginName
}

func (r *foldingRenderer) Use(p plugin.Plugin) {
	log := r.log.With(slog.String("plugin", p.Name()))

	if pr, ok := p.(plugin.Renderer); ok {
		r.plugins = append(r.plugins, pr)
	} else {
		m := fmt.Sprintf("failed to add plugin %q, since it doesn't implement plugin.Renderer", p.Name())
		log.Error(m)
		if r.panicOnInit {
			panic(fmt.Sprintf("%s: %s", foldingRendererPluginName, m))
		}
	}
}

func (r *foldingRenderer) Render(src fs.File, w io.Writer) error {
	if len(r.plugins) == 0 {
		_, err := io.Copy(w, src)
		return err
	}

	f, err := newFoldignFile(src)
	if err != nil {
		return err
	}

	for _, p := range r.plugins {
		err := p.Render(f, f)
		if err != nil {
			return err
		}

		if err := f.Fold(); err != nil {
			return err
		}
	}

	_, err = io.Copy(w, f)
	return err
}

type foldingFile struct {
	fs.File
	read   *bytes.Buffer
	writer *bytes.Buffer
}

func newFoldignFile(f fs.File) (*foldingFile, error) {
	var r, w bytes.Buffer

	if _, err := io.Copy(&r, f); err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	return &foldingFile{File: f, read: &r, writer: &w}, nil
}

func (f *foldingFile) Close() error {
	return nil
}

func (f *foldingFile) Read(p []byte) (int, error) {
	return f.read.Read(p)
}

func (f *foldingFile) Write(p []byte) (int, error) {
	return f.writer.Write(p)
}

func (f *foldingFile) Fold() error {
	f.read.Reset()
	if _, err := io.Copy(f.writer, f.read); err != nil {
		return err
	}
	f.writer.Reset()
	return nil
}
