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

type TypedMetadata interface {
	Metadata
	GetBool(key string) (bool, error)

	GetString(key string) (string, error)

	GetInt(key string) (int, error)
	GetInt8(key string) (int8, error)
	GetInt16(key string) (int16, error)
	GetInt32(key string) (int32, error)
	GetInt64(key string) (int64, error)

	GetUInt(key string) (uint, error)
	GetUInt8(key string) (uint8, error)
	GetUInt16(key string) (uint16, error)
	GetUInt32(key string) (uint32, error)
	GetUInt64(key string) (uint64, error)
	GetUIntPtr(key string) (uintptr, error)

	GetByte(key string) (byte, error)

	GetRune(key string) (rune, error)

	GetFloat32(key string) (float32, error)
	GetFloat64(key string) (float64, error)

	GetComplex64(key string) (complex64, error)
	GetComplex128(key string) (complex128, error)
}

