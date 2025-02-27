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
	"fmt"
	"net/http"
	"strings"
)

type Route interface {
	http.Handler
	fmt.Stringer
}

type RouteGroup interface {
	Routes() []Route
}

type route struct {
	method string
	host   string
	path   string

	handler http.Handler
}

func newRoute(method, host, path string, handler http.Handler) Route {
	return &route{method, host, path, handler}
}

func (r *route) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.handler.ServeHTTP(w, req)
}

func (r *route) String() string {
	path := r.host + r.path

	if r.method != "" {
		path = r.method + " " + path
	}

	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	if strings.HasSuffix(path, "...}/") {
		path = strings.TrimSuffix(path, "/")
	}

	return path
}
