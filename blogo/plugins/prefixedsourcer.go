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
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"strings"

	"forge.capytal.company/loreddev/x/blogo/metadata"
	"forge.capytal.company/loreddev/x/blogo/plugin"
)

const prefixedSourcerName = "blogo-prefixedsourcer-sourcer"

type PrefixedSourcer interface {
	plugin.Sourcer
	plugin.WithPlugins
	UseNamed(string, plugin.Plugin)
}

type prefixedSourcer struct {
	sources map[string]plugin.Sourcer

	prefixSeparator  string
	acceptDuplicated bool

	panicOnInit       bool
	skipOnSourceError bool
	skipOnFSError     bool

	log *slog.Logger
}

type PrefixedSourcerOpts struct {
	PrefixSeparator  string
	AcceptDuplicated bool

	NotPanicOnInit       bool
	NotSkipOnHexError    bool
	NotSkipOnSourceError bool
	NotSkipOnFSError     bool

	Logger *slog.Logger
}

func NewPrefixedSourcer(opts ...PrefixedSourcerOpts) PrefixedSourcer {
	opt := PrefixedSourcerOpts{}
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.PrefixSeparator == "" {
		opt.PrefixSeparator = "/"
	}

	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	opt.Logger = opt.Logger.WithGroup(prefixedSourcerName)

	return &prefixedSourcer{
		sources: map[string]plugin.Sourcer{},

		prefixSeparator:  opt.PrefixSeparator,
		acceptDuplicated: opt.AcceptDuplicated,

		panicOnInit:       !opt.NotPanicOnInit,
		skipOnSourceError: !opt.NotSkipOnSourceError,
		skipOnFSError:     !opt.NotSkipOnFSError,

		log: opt.Logger,
	}
}

func (s *prefixedSourcer) Name() string {
	return prefixedSourcerName
}

func (s *prefixedSourcer) Use(plugin plugin.Plugin) {
	s.UseNamed(plugin.Name(), plugin)
}

func (s *prefixedSourcer) UseNamed(prefix string, p plugin.Plugin) {
	log := s.log.With(slog.String("plugin", p.Name()), slog.String("prefix", prefix))

	var sourcer plugin.Sourcer
	if ps, ok := p.(plugin.Sourcer); ok {
		sourcer = ps
	} else {
		m := fmt.Sprintf("failed to add plugin %q (with prefix %q), since it doesn't implement SourcerPlugin", p.Name(), prefix)
		log.Error(m)
		if s.panicOnInit {
			panic(fmt.Sprintf("%s: %s", prefixedSourcerName, m))
		}
	}

	if _, ok := s.sources[prefix]; ok && !s.acceptDuplicated {
		m := fmt.Sprintf(
			"duplicated prefix (%q) for plugin %q",
			prefix,
			p.Name(),
		)
		log.Error(m)
		if s.panicOnInit {
			panic(fmt.Sprintf("%s: %s", prefixedSourcerName, m))
		}
		return
	}

	log.Debug(fmt.Sprintf("Added sourcer plugin, with prefix %q", prefix))
	s.sources[prefix] = sourcer
}

func (s *prefixedSourcer) Source() (fs.FS, error) {
	log := s.log

	fileSystems := make(map[string]fs.FS, len(s.sources))

	for a, ps := range s.sources {
		log = log.With(slog.String("plugin", ps.Name()), slog.String("prefix", a))
		log.Info("Sourcing file system of plugin")

		f, err := ps.Source()
		if err != nil && s.skipOnSourceError {
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
		prefixSeparator: s.prefixSeparator,
	}, nil
}

type prefixedSourcerFS struct {
	fileSystems     map[string]fs.FS
	prefixSeparator string
}

func (pf *prefixedSourcerFS) Metadata() metadata.Metadata {
	ms := []metadata.Metadata{}
	for _, v := range pf.fileSystems {
		if m, err := metadata.GetMetadata(v); err == nil {
			ms = append(ms, m)
		}
	}
	return metadata.Join(ms...)
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
