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
	"cmp"
	"slices"
)

const priorityGroupPluginName = "blogo-prioritygroup-group"

type priorityGroup struct {
	plugins []Plugin
}

type PriorityGroup interface {
	PluginWithPlugins
}

func NewPriorityGroup(plugins ...Plugin) PriorityGroup {
	return &priorityGroup{plugins}
}

func (p *priorityGroup) Name() string {
	return priorityGroupPluginName
}

func (p *priorityGroup) Use(plugin Plugin) {
	p.plugins = append(p.plugins, plugin)
}

func (p *priorityGroup) Plugins() []Plugin {
	slices.SortStableFunc(p.plugins, func(a Plugin, b Plugin) int {
		return cmp.Compare(p.getPriority(a, b), p.getPriority(b, a))
	})
	return p.plugins
}

func (p *priorityGroup) getPriority(plugin Plugin, cmp Plugin) int {
	if plg, ok := plugin.(PluginWithDynamicPriority); ok {
		return plg.Priority(cmp)
	} else if plg, ok := plugin.(PluginWithPriority); ok {
		return plg.Priority()
	} else {
		return 0
	}
}

type PluginWithPriority interface {
	Plugin
	Priority() int
}

type PluginWithDynamicPriority interface {
	Plugin
	Priority(Plugin) int
}
