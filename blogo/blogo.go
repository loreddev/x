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
	"io"
	"log/slog"
	"net/http"

	"forge.capytal.company/loreddev/x/blogo/core"
	"forge.capytal.company/loreddev/x/blogo/plugin"
	"forge.capytal.company/loreddev/x/blogo/plugins"
	"forge.capytal.company/loreddev/x/tinyssert"
)

func New(opts ...Opts) Blogo {
	opt := Opts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Assertions == nil {
		opt.Assertions = tinyssert.NewDisabledAssertions()
	}
	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}

	if opt.FallbackRenderer == nil {
		opt.FallbackRenderer = plugins.NewPlainText()
	}
	if opt.MultiRenderer == nil {
		opt.MultiRenderer = plugins.NewMultiRenderer()
	}
	if opt.FallbackSourcer == nil {
		opt.FallbackSourcer = plugins.NewEmptySourcer()
	}
	if opt.MultiSourcer == nil {
		opt.MultiSourcer = plugins.NewMultiSourcer()
	}

	return &blogo{
		plugins: []plugin.Plugin{},

		fallbackRenderer: opt.FallbackRenderer,
		multiRenderer:    opt.MultiRenderer,
		fallbackSourcer:  opt.FallbackSourcer,
		multiSourcer:     opt.MultiSourcer,

		assert: opt.Assertions,
		log:    opt.Logger,
	}
}

type Blogo interface {
	Use(plugin.Plugin)
	Init()
	http.Handler
}

type Opts struct {
	FallbackRenderer plugin.Renderer
	MultiRenderer    interface {
		plugin.Renderer
		plugin.WithPlugins
	}
	FallbackSourcer plugin.Sourcer
	MultiSourcer    interface {
		plugin.Sourcer
		plugin.WithPlugins
	}

	Assertions tinyssert.Assertions
	Logger     *slog.Logger
}

type blogo struct {
	plugins []plugin.Plugin

	fallbackRenderer plugin.Renderer
	multiRenderer    interface {
		plugin.Renderer
		plugin.WithPlugins
	}
	fallbackSourcer plugin.Sourcer
	multiSourcer    interface {
		plugin.Sourcer
		plugin.WithPlugins
	}

	server http.Handler

	assert tinyssert.Assertions
	log    *slog.Logger
}

func (b *blogo) Use(p plugin.Plugin) {
	b.assert.NotNil(p, "Plugin definition should not be nil")
	b.assert.NotNil(b.plugins, "Plugins needs to be not-nil")
	b.assert.NotNil(b.log)

	log := b.log.With(slog.String("plugin", p.Name()))

	if p, ok := p.(plugin.Group); ok {
		log.Debug("Plugin group found, adding it's plugins")
		for _, p := range p.Plugins() {
			b.Use(p)
		}
	}

	b.plugins = append(b.plugins, p)
}

func (b *blogo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.assert.NotNil(b.log)
	b.assert.NotNil(w)
	b.assert.NotNil(r)

	if b.server != nil {
		b.server.ServeHTTP(w, r)
		return
	}

	log := b.log.With()
	log.Debug("Core server not initialized")

	b.Init()

	b.server.ServeHTTP(w, r)
}

func (b *blogo) Init() {
	b.assert.NotNil(b.plugins, "Plugins needs to be not-nil")
	b.assert.NotNil(b.log)

	log := b.log.With()
	log.Debug("Initializing Blogo plugins")

	sourcer := b.initSourcer()
	renderer := b.initRenderer()

	log.Debug("Constructing Blogo server")

	b.server = core.NewServer(sourcer, renderer, core.ServerOpts{
		Assertions: b.assert,
		Logger:     b.log.WithGroup("server"),
	})

	log.Debug("Server constructed")
}

func (b *blogo) initRenderer() plugin.Renderer {
	b.assert.NotNil(b.plugins, "Plugins needs to be not-nil")
	b.assert.NotNil(b.fallbackRenderer, "FallbackRenderer needs to be not-nil")
	b.assert.NotNil(b.multiRenderer, "MultiRenderer needs to be not-nil")
	b.assert.NotNil(b.log)

	log := b.log.With()
	log.Debug("Initializing Blogo Renderer plugins")

	renderers := []plugin.Renderer{}

	for _, p := range b.plugins {
		if r, ok := p.(plugin.Renderer); ok {
			log.Debug("Adding Renderer", slog.String("sourcer", r.Name()))

			renderers = append(renderers, r)
		}
	}

	if len(renderers) == 0 {
		log.Debug("No Renderer avaiable, using %q as fallback",
			slog.String("renderer", b.fallbackRenderer.Name()))

		return b.fallbackRenderer
	}

	if len(renderers) == 1 {
		log.Debug("Just one Renderer found, using it directly",
			slog.String("renderer", renderers[0].Name()))

		return renderers[0]
	}

	log.Debug("Multiple Renderers found, using MultiRenderer to combine them",
		slog.String("renderer", b.multiRenderer.Name()),
	)
	for _, r := range renderers {
		b.multiRenderer.Use(r)
	}

	return b.multiRenderer
}

func (b *blogo) initSourcer() plugin.Sourcer {
	b.assert.NotNil(b.plugins, "Plugins needs to be not-nil")
	b.assert.NotNil(b.fallbackSourcer, "FallbackSourcer needs to be not-nil")
	b.assert.NotNil(b.multiSourcer, "MultiSourcer needs to be not-nil")
	b.assert.NotNil(b.log)

	log := b.log.With()
	log.Debug("Initializing Blogo Sourcer plugins")

	sourcers := []plugin.Sourcer{}

	for _, p := range b.plugins {
		if s, ok := p.(plugin.Sourcer); ok {
			log.Debug("Adding Sourcer", slog.String("sourcer", s.Name()))

			sourcers = append(sourcers, s)
		}
	}

	if len(sourcers) == 0 {
		log.Debug("No Sourcer avaiable, using %q as fallback",
			slog.String("sourcer", b.fallbackSourcer.Name()))

		return b.fallbackSourcer
	}

	if len(sourcers) == 1 {
		log.Debug("Just one Sourcer found, using it directly",
			slog.String("sourcer", sourcers[0].Name()))

		return sourcers[0]
	}

	log.Debug("Multiple Sourcers found, using MultiSourcer to combine them",
		slog.String("sourcer", b.multiSourcer.Name()),
	)
	for _, s := range sourcers {
		b.multiSourcer.Use(s)
	}

	return b.multiSourcer
}
