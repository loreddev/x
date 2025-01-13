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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
)

var ErrRendererNotSupportedFile = errors.New("this file is not supported by renderer")

const multiRendererPluginName = "blogo-multirenderer-renderer"

type MultiRenderer interface {
	RendererPlugin
	Use(Plugin)
}

type multiRenderer struct {
	renderers []RendererPlugin

	skipOnError bool
	panicOnInit bool

	log *slog.Logger
}

type MultiRendererOpts struct {
	NotSkipOnError bool
	NotPanicOnInit bool
	Logger         *slog.Logger
}

func NewMultiRenderer(opts ...MultiRendererOpts) MultiRenderer {
	opt := MultiRendererOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	opt.Logger = opt.Logger.WithGroup(multiRendererPluginName)

	return &multiRenderer{
		renderers: []RendererPlugin{},

		skipOnError: !opt.NotSkipOnError,
		panicOnInit: !opt.NotPanicOnInit,

		log: opt.Logger,
	}
}

func (p *multiRenderer) Name() string {
	return multiRendererPluginName
}

func (p *multiRenderer) Use(plugin Plugin) {
	log := p.log.With(slog.String("plugin", plugin.Name()))

	if plg, ok := plugin.(RendererPlugin); ok {
		log.Debug("Added renderer plugin")
		p.renderers = append(p.renderers, plg)
	} else {
		m := fmt.Sprintf("failed to add plugin %q, since it doesn't implement RendererPlugin", plugin.Name())
		log.Error(m)
		if p.panicOnInit {
			panic(fmt.Sprintf("%s: %s", p.Name(), m))
		}
	}
}

func (p *multiRenderer) Render(f fs.File, w io.Writer) error {
	mf := newMultiRendererFile(f)
	for _, r := range p.renderers {
		log := p.log.With(slog.String("plugin", r.Name()))

		log.Debug("Trying to render with plugin")
		err := r.Render(f, w)

		if err == nil {
			break
		}

		if !p.skipOnError && !errors.Is(err, ErrRendererNotSupportedFile) {
			log.Error("Failed to render using plugin", slog.String("error", err.Error()))
			return errors.Join(fmt.Errorf("failed to render using plugin %q", p.Name()), err)
		}

		log.Debug("Unable to render using plugin", slog.String("error", err.Error()))
		log.Debug("Resetting file for next read")

		if err := mf.Reset(); err != nil {
			log.Error("Failed to reset file read offset", slog.String("error", err.Error()))
			return errors.Join(fmt.Errorf("failed to reset file read offset"), err)
		}
	}

	return nil
}

type multiRendererFile struct {
	fs.File
	buf    *bytes.Buffer
	reader io.Reader
}

func newMultiRendererFile(f fs.File) *multiRendererFile {
	if _, ok := f.(io.Seeker); ok {
		return &multiRendererFile{
			File:   f,
			reader: f,
		}
	}

	var buf bytes.Buffer
	return &multiRendererFile{
		File:   f,
		reader: io.TeeReader(f, &buf),
		buf:    &buf,
	}
}

func (f *multiRendererFile) Read(p []byte) (int, error) {
	return f.reader.Read(p)
}

func (f *multiRendererFile) Reset() error {
	if s, ok := f.File.(io.Seeker); ok {
		_, err := s.Seek(0, io.SeekStart)
		return err
	}
	var buf bytes.Buffer
	r := io.MultiReader(f.buf, f.File)

	f.reader = io.TeeReader(r, &buf)
	f.buf = &buf

	return nil
}
