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
	"bytes"
	"fmt"
	"io"
	"log/slog"
)

const foldingRendererPluginName = "blogo-foldingrenderer-renderer"

type foldingRenderer struct {
	plugins []RendererPlugin

	panicOnInit bool

	log *slog.Logger
}

type FoldingRendererOpts struct {
	PanicOnInit bool

	Logger *slog.Logger
}

type FoldingRenderer interface {
	PluginWithPlugins
	RendererPlugin
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
		plugins: []RendererPlugin{},

		log: opt.Logger,
	}
}

func (p *foldingRenderer) Name() string {
	return foldingRendererPluginName
}

func (p *foldingRenderer) Use(plugin Plugin) {
	log := p.log.With(slog.String("plugin", plugin.Name()))

	if plg, ok := plugin.(RendererPlugin); ok {
		p.plugins = append(p.plugins, plg)
	} else {
		m := fmt.Sprintf("failed to add plugin %q, since it doesn't implement RendererPlugin", plugin.Name())
		log.Error(m)
		if p.panicOnInit {
			panic(fmt.Sprintf("%s: %s", multiRendererPluginName, m))
		}
	}
}

func (p *foldingRenderer) Render(src File, w io.Writer) error {
	if len(p.plugins) == 0 {
		_, err := io.Copy(w, src)
		return err
	}

	f, err := newFoldignFile(src)
	if err != nil {
		return err
	}

	for _, r := range p.plugins {
		err := r.Render(f, f)
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
	File
	read   *bytes.Buffer
	writer *bytes.Buffer
}

func newFoldignFile(f File) (*foldingFile, error) {
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
