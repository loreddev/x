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
)

const multiSourcerPluginName = "blogo-multisourcer-sourcer"

type MultiSourcer interface {
	SourcerPlugin
	Use(Plugin)
}

type multiSourcer struct {
	sources []SourcerPlugin

	panicOnInit       bool
	skipOnSourceError bool
	skipOnFSError     bool

	log *slog.Logger
}

type MultiSourcerOpts struct {
	NotPanicOnInit       bool
	NotSkipOnSourceError bool
	NotSkipOnFSError     bool

	Logger *slog.Logger
}

func NewMultiSourcer(opts ...MultiSourcerOpts) MultiSourcer {
	opt := MultiSourcerOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	opt.Logger = opt.Logger.WithGroup(multiSourcerPluginName)

	return &multiSourcer{
		sources: []SourcerPlugin{},

		panicOnInit:       !opt.NotPanicOnInit,
		skipOnSourceError: !opt.NotSkipOnSourceError,
		skipOnFSError:     !opt.NotSkipOnFSError,

		log: opt.Logger,
	}
}

func (p *multiSourcer) Name() string {
	return multiSourcerPluginName
}

func (p *multiSourcer) Use(plugin Plugin) {
	log := p.log.With(slog.String("plugin", plugin.Name()))

	if plg, ok := plugin.(SourcerPlugin); ok {
		log.Debug("Added renderer plugin")
		p.sources = append(p.sources, plg)
	} else {
		m := fmt.Sprintf("failed to add plugin %q, since it doesn't implement SourcerPlugin", plugin.Name())
		log.Error(m)
		if p.panicOnInit {
			panic(fmt.Sprintf("%s: %s", multiRendererPluginName, m))
		}
	}
}

func (p *multiSourcer) Source() (fs.FS, error) {
	log := p.log

	fileSystems := []fs.FS{}

	for _, s := range p.sources {
		log = log.With(slog.String("plugin", p.Name()))
		log.Info("Sourcing file system of plugin")

		f, err := s.Source()
		if err != nil && p.skipOnSourceError {
			log.Error(
				"Failed to source file system of plugin, skipping",
				slog.String("error", err.Error()),
			)
		} else if err != nil {
			log.Error(
				"Failed to source file system of plugin, returning error",
				slog.String("error", err.Error()),
			)
			return f, err
		}

		fileSystems = append(fileSystems, f)
	}

	f := make([]fs.FS, len(fileSystems), len(fileSystems))
	for i := range f {
		f[i] = fileSystems[i]
	}

	return &multiSourcerFS{
		fileSystems: f,
		skipOnError: p.skipOnFSError,
	}, nil
}

type multiSourcerFS struct {
	fileSystems []fs.FS
	skipOnError bool
}

func (mf *multiSourcerFS) Open(name string) (fs.File, error) {
	for _, f := range mf.fileSystems {
		file, err := f.Open(name)

		if err != nil && !errors.Is(err, fs.ErrNotExist) && !mf.skipOnError {
			return file, err
		}

		if err == nil {
			return file, err
		}
	}

	return nil, fs.ErrNotExist
}
