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

import "errors"

var (
	ErrMetadataNotFound      = errors.New("key in metadata was not found")
	ErrMetadataIncorrectType = errors.New("key in metadata is not of the provided type")
	ErrMetadataImmutable     = errors.New("key in metadata cannot be set to another value")
	ErrMetadataNotEmpty      = errors.New("key in metadata is not empty")
)

type Metadata interface {
	Get(key string) (any, error)
	Set(key string, v any, strict ...bool) error
	Delete(key string, strict ...bool) error
}

type metadataMap map[string]any

func MetadataMap(m map[string]any) Metadata {
	if m == nil {
		m = map[string]any{}
	}
	return metadataMap(m)
}

func (m metadataMap) Get(key string) (any, error) {
	v, ok := m[key]
	if !ok {
		return nil, ErrMetadataNotFound
	}
	return v, nil
}

func (m metadataMap) Set(key string, v any, strict ...bool) error {
	if _, ok := m[key]; ok && len(strict) > 0 && strict[0] {
		return ErrMetadataNotEmpty
	}
	m[key] = v
	return nil
}

func (m metadataMap) Delete(key string, strict ...bool) error {
	if _, ok := m[key]; ok && len(strict) > 0 && strict[0] {
		return ErrMetadataNotEmpty
	}
	delete(m, key)
	return nil
}

type multiFSMetadata struct {
	Metadata
	fileSystems []FS
}

func NewMultiFSMetadata(fileSytems []FS) Metadata {
	return &multiFSMetadata{
		Metadata:    MetadataMap(map[string]any{}),
		fileSystems: fileSytems,
	}
}

func (m *multiFSMetadata) Get(key string) (any, error) {
	if v, err := m.Metadata.Get(key); err == nil {
		return v, nil
	}

	for _, m := range m.fileSystems {
		v, err := m.Metadata().Get(key)
		if err == nil {
			return v, nil
		}
	}
	return nil, ErrMetadataNotFound
}

