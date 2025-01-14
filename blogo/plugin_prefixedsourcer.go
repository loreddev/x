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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"strings"
)

const prefixedSourcerPluginName = "blogo-prefixedsourcer-sourcer"

type PrefixedSourcer interface {
	SourcerPlugin
	PluginWithPlugins
	UseNamed(string, Plugin)
}

type prefixedSourcer struct {
	sources map[string]SourcerPlugin

	prefixSeparator  string
	acceptDuplicated bool

	panicOnInit       bool
	skipOnHexError    bool
	skipOnSourceError bool
	skipOnFSError     bool

	hexFunc HexFunc

	log *slog.Logger
}

type PrefixedSourcerOpts struct {
	PrefixSeparator  string
	AcceptDuplicated bool

	NotPanicOnInit       bool
	NotSkipOnHexError    bool
	NotSkipOnSourceError bool
	NotSkipOnFSError     bool

	HexFunc HexFunc

	Logger *slog.Logger
}

type HexFunc = func(n int) (string, error)

func hexFunc(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func NewPrefixedSourcer(opts ...PrefixedSourcerOpts) PrefixedSourcer {
	opt := PrefixedSourcerOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.PrefixSeparator == "" {
		opt.PrefixSeparator = "/"
	}

	if opt.HexFunc == nil {
		opt.HexFunc = hexFunc
	}

	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	opt.Logger = opt.Logger.WithGroup(prefixedSourcerPluginName)

	return &prefixedSourcer{
		sources: map[string]SourcerPlugin{},

		prefixSeparator:  opt.PrefixSeparator,
		acceptDuplicated: opt.AcceptDuplicated,

		panicOnInit:       !opt.NotPanicOnInit,
		skipOnSourceError: !opt.NotSkipOnSourceError,
		skipOnFSError:     !opt.NotSkipOnFSError,

		log: opt.Logger,
	}
}

func (p *prefixedSourcer) Name() string {
	return prefixedSourcerPluginName
}

func (p *prefixedSourcer) Use(plugin Plugin) {
	p.UseNamed(plugin.Name(), plugin)
}

func (p *prefixedSourcer) UseNamed(prefix string, plugin Plugin) {
	log := p.log.With(slog.String("plugin", plugin.Name()), slog.String("prefix", prefix))

	var sourcer SourcerPlugin
	if plg, ok := plugin.(SourcerPlugin); ok {
		sourcer = plg
	} else {
		m := fmt.Sprintf("failed to add plugin %q (with prefix %q), since it doesn't implement SourcerPlugin", plugin.Name(), prefix)
		log.Error(m)
		if p.panicOnInit {
			panic(fmt.Sprintf("%s: %s", multiRendererPluginName, m))
		}
	}

	if _, ok := p.sources[prefix]; ok && !p.acceptDuplicated {
		m := fmt.Sprintf(
			"duplicated prefix (%q) for plugin %q",
			prefix,
			plugin.Name(),
		)
		log.Error(m)
		if p.panicOnInit {
			panic(fmt.Sprintf("%s: %s", multiRendererPluginName, m))
		}
		return
	}

	log.Debug(fmt.Sprintf("Added sourcer plugin, with prefix %q", prefix))
	p.sources[prefix] = sourcer
}

func (p *prefixedSourcer) Source() (fs.FS, error) {
	log := p.log

	fileSystems := make(map[string]fs.FS, len(p.sources))

	for a, s := range p.sources {
		log = log.With(slog.String("plugin", p.Name()), slog.String("prefix", a))
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

		fileSystems[a] = f
	}

	return &prefixedSourcerFS{
		fileSystems:     fileSystems,
		prefixSeparator: p.prefixSeparator,
	}, nil
}

type prefixedSourcerFS struct {
	fileSystems     map[string]fs.FS
	prefixSeparator string
}

func (pf *prefixedSourcerFS) Open(name string) (fs.File, error) {
	prefix, path, found := strings.Cut(name, pf.prefixSeparator)
	if !found {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	if f, ok := pf.fileSystems[prefix]; ok {
		return f.Open(path)
	}

	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
}
