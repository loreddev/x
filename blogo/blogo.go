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
}

func (b *Blogo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := b.log.With(slog.String("path", r.URL.Path))

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

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}

		if n == 0 {
			break
		}

		_, err = w.Write(buf[:n])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(err.Error()))
		}
	}

	b.log.Debug("Finished responding file")
}

func (b *Blogo) Init() error {
	b.log.Debug("Initializing blogo")

	if len(b.sources) == 0 {
		b.log.Debug("No SourcerPlugin found, using default one")
		b.Use(&defaultSourcer{})
	}

	fs, err := b.sources[0].Source() // TOOD: Support for multiple sources (via another plugin or built-in, with prefixes or not)
	if err != nil {
		return errors.Join(errors.New("failed to source files"), err)
	}
	b.files = fs

	return nil
}
