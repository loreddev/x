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

package smalltrip

import (
	"log/slog"
	"net/http"

	"forge.capytal.company/loreddev/x/tinyssert"
)

type Option func(*router)

func WithAssertions(assertions tinyssert.Assertions) Option {
	return func(r *router) {
		r.assert = assertions
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(r *router) {
		r.log = logger
	}
}

func WithServeMux(mux *http.ServeMux) Option {
	return func(r *router) {
		r.mux = mux
	}
}
