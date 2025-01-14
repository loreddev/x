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

