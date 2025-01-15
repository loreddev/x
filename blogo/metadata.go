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

type joinedMetadata struct {
	ms []Metadata
	m  Metadata
}

func JoinMetadata(ms ...Metadata) Metadata {
	jm := []Metadata{}
	for _, m := range ms {
		if ms, ok := m.(*joinedMetadata); ok {
			jm = append(jm, ms.m)
			jm = append(jm, ms.ms...)
		} else if m != nil {
			jm = append(jm, m)
		}
	}
	return &joinedMetadata{
		ms: jm,
		m:  MetadataMap(map[string]any{}),
	}
}

func (jm *joinedMetadata) Get(key string) (any, error) {
	if v, err := jm.m.Get(key); err == nil {
		return v, nil
	}
	for _, m := range jm.ms {
		v, err := m.Get(key)
		if err == nil {
			return v, nil
		}
	}
	return nil, ErrMetadataNotEmpty
}

func (jm *joinedMetadata) Set(key string, v any, strict ...bool) error {
	if _, err := jm.m.Get(key); err == nil {
		return jm.m.Set(key, v, strict...)
	}

	for _, m := range jm.ms {
		_, err := m.Get(key)
		if err == nil {
			return m.Set(key, v, strict...)
		} else if errors.Is(err, ErrMetadataImmutable) {
			return err
		}
	}

	return jm.m.Set(key, v, strict...)
}

func (jm *joinedMetadata) Delete(key string, strict ...bool) error {
	if _, err := jm.m.Get(key); err == nil {
		return jm.m.Delete(key, strict...)
	}

	for _, m := range jm.ms {
		_, err := m.Get(key)
		if err == nil {
			return m.Delete(key, strict...)
		} else if errors.Is(err, ErrMetadataImmutable) {
			return err
		}
	}

	return jm.m.Delete(key, strict...)
}

type immutableMetadata struct {
	Metadata
}

func ImmutableMetadata(m Metadata) Metadata {
	return &immutableMetadata{m}
}

func (m *immutableMetadata) Set(key string, v any, strict ...bool) error {
	return ErrMetadataImmutable
}

func (m *immutableMetadata) Delete(key string, strict ...bool) error {
	return ErrMetadataImmutable
}

type TypedMetadata struct {
	Metadata
}

func NewTypedMetadata(m Metadata) *TypedMetadata {
	return &TypedMetadata{m}
}

func (m *TypedMetadata) GetBool(key string) (bool, error) {
	return GetTyped[bool](m, key)
}

func (m *TypedMetadata) GetString(key string) (string, error) {
	return GetTyped[string](m, key)
}

func (m *TypedMetadata) GetInt(key string) (int, error) {
	return GetTyped[int](m, key)
}

func (m *TypedMetadata) GetInt8(key string) (int8, error) {
	return GetTyped[int8](m, key)
}

func (m *TypedMetadata) GetInt16(key string) (int16, error) {
	return GetTyped[int16](m, key)
}

func (m *TypedMetadata) GetInt32(key string) (int32, error) {
	return GetTyped[int32](m, key)
}

func (m *TypedMetadata) GetInt64(key string) (int64, error) {
	return GetTyped[int64](m, key)
}

func (m *TypedMetadata) GetUint(key string) (uint, error) {
	return GetTyped[uint](m, key)
}

func (m *TypedMetadata) GetUint8(key string) (uint8, error) {
	return GetTyped[uint8](m, key)
}

func (m *TypedMetadata) GetUint16(key string) (uint16, error) {
	return GetTyped[uint16](m, key)
}

func (m *TypedMetadata) GetUint32(key string) (uint32, error) {
	return GetTyped[uint32](m, key)
}

func (m *TypedMetadata) GetUint64(key string) (uint64, error) {
	return GetTyped[uint64](m, key)
}

func (m *TypedMetadata) GetByte(key string) (byte, error) {
	return GetTyped[byte](m, key)
}

func (m *TypedMetadata) GetRune(key string) (rune, error) {
	return GetTyped[rune](m, key)
}

func (m *TypedMetadata) GetFloat32(key string) (float32, error) {
	return GetTyped[float32](m, key)
}

func (m *TypedMetadata) GetFloat64(key string) (float64, error) {
	return GetTyped[float64](m, key)
}

func (m *TypedMetadata) GetComplex64(key string) (complex64, error) {
	return GetTyped[complex64](m, key)
}

func (m *TypedMetadata) GetComplex128(key string) (complex128, error) {
	return GetTyped[complex128](m, key)
}

func GetTyped[T any](m Metadata, key string) (T, error) {
	var z T

	v, err := m.Get(key)
	if err != nil {
		return z, err
	}

	if v, ok := v.(T); ok {
		return v, nil
	} else {
		return z, ErrMetadataIncorrectType
	}
}
