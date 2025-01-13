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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
)

type Options struct {
	Logger *slog.Logger
}

type Blogo struct {
	files fs.FS

	sources   []SourcerPlugin
	renderers []RendererPlugin

	log   *slog.Logger
	panic bool
}

func New(opts ...Options) *Blogo {
	opt := Options{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	} else {
		opt.Logger = opt.Logger.WithGroup("blogo")
	}

	return &Blogo{
		files:   nil,
		sources: []SourcerPlugin{},
		log:     opt.Logger,
		panic:   true, // TODO
	}
}

func (b *Blogo) Use(p Plugin) {
	log := b.log.With(slog.String("plugin", p.Name()))

	if p, ok := p.(SourcerPlugin); ok {
		log.Debug("Added plugin", slog.String("type", "SourcerPlugin"))
		b.sources = append(b.sources, p)
	}
	if p, ok := p.(RendererPlugin); ok {
		log.Debug("Added plugin", slog.String("type", "RenderPlugin"))
		b.renderers = append(b.renderers, p)
	}
}

func (b *Blogo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := b.log.With(slog.String("step", "SERVE"), slog.String("path", r.URL.Path))

	log.Debug("Serving endpoint")

	if b.files == nil {
		log.Debug("No files in Blogo engine, initializing files")

		err := b.Init()
		if err != nil {
			log.Error("Failed to initialize files")

			err = errors.Join(errors.New("failed to initialize Blogo engine on first request"), err)
			if b.panic {
				panic(err.Error())
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(err.Error()))
			}
			return
		}
	}

	path := strings.Trim(r.URL.Path, "/")
	if path == "" || path == "/" {
		path = "."
	}

	f, err := b.files.Open(path)

	if errors.Is(err, fs.ErrNotExist) {
		log.Error("Failed to read file", slog.String("error", err.Error()))

		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		log.Error("Failed to read file", slog.String("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()

	log.Debug("Writing response file")

	log.Debug("Rendering file")

	err = b.render(f, w)
	if err != nil {
		log.Error("Failed to render file", slog.String("error", err.Error()))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	log.Debug("Finished responding file")
}

func (b *Blogo) Init() error {
	log := b.log.With(slog.String("step", "INITIALIZATION"))
	log.Debug("Initializing blogo")

	if len(b.sources) == 0 {
		log.Debug("No SourcerPlugin found, using default one")
		b.Use(&defaultSourcer{})
	}

	if len(b.renderers) == 0 {
		renderer := NewPlainTextRenderer()
		log.Debug(
			fmt.Sprintf(
				"No RendererPlugin plugin found, adding %q as fallback renderer",
				renderer.Name(),
			),
		)
		b.Use(renderer)
	}

	fs, err := b.sources[0].Source() // TODO: Support for multiple sources (via another plugin or built-in, with prefixes or not)
	if err != nil {
		return errors.Join(errors.New("failed to source files"), err)
	}
	b.files = fs

	return nil
}

func (b *Blogo) render(src fs.File, w io.Writer) error {
	log := b.log.With(slog.String("step", "RENDERING"))

	if len(b.renderers) == 1 {
		log.Debug(
			"Just one renderer found, using it directly",
			slog.String("plugin", b.renderers[0].Name()),
		)
		return b.renderers[0].Render(src, w)
	}

	log.Debug("Multiple renderers found, initializing built-in multi-renderer plugin")

	multi := NewMultiRenderer(MultiRendererOpts{
		NotSkipOnError: false,
		NotPanicOnInit: true,
		Logger:         log,
	})

	for _, r := range b.renderers {
		log.Debug("Adding plugin to multi-renderer", slog.String("plugin", r.Name()))
		multi.Use(r)
	}

	log.Debug("Overriding renderers slice")

	b.renderers = make([]RendererPlugin, 1)
	b.renderers[0] = multi

	return b.render(src, w)
}
