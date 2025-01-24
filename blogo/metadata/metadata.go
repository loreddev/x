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

package metadata

import (
	"errors"
	"reflect"
)

var (
	ErrImmutable   = errors.New("metadata is immutable")
	ErrInvalidType = errors.New("key is not of specified type")
	ErrNotFound    = errors.New("value not found metadata")
	ErrNoMetadata  = errors.New("type does not implement or has metadata")
)

func GetTyped[T any](m any, key string) (T, error) {
	var z T

	if m, ok := m.(TypedMetadata); ok {
		v, err := getTypedFromTyped[T](m, key)
		if v, ok := v.(T); ok && err == nil {
			return v, err
		}
	}

	v, err := Get(m, key)
	if err != nil {
		return z, err
	}

	if v, ok := v.(T); ok {
		return v, nil
	}

	vv, zv := reflect.ValueOf(v), reflect.ValueOf(z)
	vt, zt := vv.Type(), zv.Type()
	if vt.ConvertibleTo(zt) {
		v = vv.Convert(zt).Interface()
		if v, ok := v.(T); ok {
			return v, nil
		}
	}

	return z, ErrInvalidType
}

func getTypedFromTyped[T any](m TypedMetadata, key string) (any, error) {
	var z T

	t := reflect.TypeOf(z)
	if t == nil {
		return z, ErrInvalidType
	}

	switch t.Kind() {
	case reflect.Bool:
		return m.GetBool(key)

	case reflect.String:
		return m.GetString(key)

	case reflect.Int:
		return m.GetInt(key)
	case reflect.Int8:
		return m.GetInt8(key)
	case reflect.Int16:
		return m.GetInt16(key)
	case reflect.Int32:
		return m.GetInt32(key)
	case reflect.Int64:
		return m.GetInt64(key)

	case reflect.Uint:
		return m.GetInt(key)
	case reflect.Uint8:
		return m.GetUint8(key)
	case reflect.Uint16:
		return m.GetUint16(key)
	case reflect.Uint32:
		return m.GetUint32(key)
	case reflect.Uint64:
		return m.GetUint64(key)
	case reflect.Uintptr:
		return m.GetUintptr(key)

	case reflect.Float32:
		return m.GetFloat32(key)
	case reflect.Float64:
		return m.GetFloat64(key)

	case reflect.Complex64:
		return m.GetComplex64(key)
	case reflect.Complex128:
		return m.GetComplex128(key)

	default:
		return m.Get(key)
	}
}

func Get(m any, key string) (any, error) {
	data, err := GetMetadata(m)
	if err != nil {
		return nil, err
	}
	return data.Get(key)
}

func Set(m any, key string, v any) error {
	data, err := GetMetadata(m)
	if err != nil {
		return err
	}
	return data.Set(key, v)
}

func Delete(m any, key string) error {
	data, err := GetMetadata(m)
	if err != nil {
		return err
	}
	return data.Delete(key)
}

func GetMetadata(m any) (Metadata, error) {
	var data Metadata

	if mt, ok := m.(Metadata); ok {
		data = mt
	} else if mfile, ok := m.(WithMetadata); ok {
		data = mfile.Metadata()
	} else {
		return nil, ErrNoMetadata
	}

	return data, nil
}

type WithMetadata interface {
	Metadata() Metadata
}

type Metadata interface {
	Get(key string) (any, error)
	Set(key string, v any) error
	Delete(key string) error
}

type Map map[string]any

func (m Map) Get(key string) (any, error) {
	if m == nil {
		return nil, ErrNotFound
	}

	if v, ok := m[key]; ok {
		return v, nil
	} else {
		return nil, ErrNotFound
	}
}

func (m Map) Set(key string, v any) error {
	if m == nil {
		return nil
	}
	m[key] = v
	return nil
}

func (m Map) Delete(key string) error {
	delete(m, key)
	return nil
}

func Join(ms ...Metadata) Metadata {
	ms = append([]Metadata{Map(make(map[string]any))}, ms...)
	return joined(ms)
}

type joined []Metadata

func (m joined) Get(key string) (any, error) {
	for _, m := range m {
		v, err := m.Get(key)
		if err == nil {
			return v, nil
		}
	}
	return nil, ErrNotFound
}

func (m joined) Set(key string, v any) error {
	for _, m := range m {
		_ = m.Set(key, v)
	}
	return nil
}

func (m joined) Delete(key string) error {
	for _, m := range m {
		_ = m.Delete(key)
	}
	return nil
}

type immutable struct{ Metadata }

func Immutable(m Map) Metadata {
	return &immutable{m}
}

func (m *immutable) Set(key string, v any) error {
	return ErrImmutable
}

func (m *immutable) Delete(key string) error {
	return ErrImmutable
}
