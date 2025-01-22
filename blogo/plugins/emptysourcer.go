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
	"forge.capytal.company/loreddev/x/blogo/fs"
	"forge.capytal.company/loreddev/x/blogo/plugin"
)

const emptySourcerPluginName = "blogo-empty-sourcer"

type emptySourcer struct{}

func NewEmptySourcer() plugin.Plugin {
	return &emptySourcer{}
}

func (p *emptySourcer) Name() string {
	return emptySourcerPluginName
}

func (p *emptySourcer) Source() (fs.FS, error) {
	return emptyFS{}, nil
}

type emptyFS struct{}

func (f emptyFS) Metadata() fs.Metadata {
	return fs.MetadataMap(map[string]any{})
}

func (f emptyFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}
