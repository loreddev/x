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

package exceptions

import (
	"errors"
	"fmt"
	"net/http"
)

type Exception struct {
	Status   int            `json:"status"`          // HTTP Status Code
	Code     string         `json:"code"`            // Application error code
	Message  string         `json:"message"`         // User friendly message
	Err      error          `json:"error,omitempty"` // Go error
	Data     map[string]any `json:"data,omitempty"`  // Additional contextual data
	Severity Severity       `json:"severity"`        // Exception level

	// Handler to be used. This is normally provided by a middleware via the
	// request context. Setting this field overrides any provided by the middleware
	// and can be used to add a handler when using a middleware is not possible.
	handler HandlerFunc `json:"-"`

	headers http.Header
}

var (
	_ fmt.Stringer = Exception{}
	_ error        = Exception{}
	_ http.Handler = Exception{}
)

func (e Exception) String() string {
	return fmt.Sprintf("%s %3d %s Exception %q", e.Severity, e.Status, e.Code, e.Message)
}

func (e Exception) Error() string {
	return e.String()
}

func (e Exception) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.handler != nil {
		e.handler(e, w, r)
	}

	e.handler = HandlerJSON(HandlerText)

	handler, ok := r.Context().Value(handlerFuncCtxKey).(HandlerFunc)
	if !ok {
		e.handler(e, w, r)
		return
	}

	handler(e, w, r)
}

func newException(options ...Option) Exception {
	e := Exception{
		Status:   http.StatusInternalServerError,
		Code:     "Internal Server Error",
		Message:  "",
		Err:      nil,
		Severity: ERROR,
	}

	for _, option := range options {
		option(&e)
	}

	return e
}

type Option = func(*Exception)

func WithStatus(s int) Option {
	return func(e *Exception) { e.Status = s }
}

func WithCode(c string) Option {
	return func(e *Exception) { e.Code = c }
}

func WithMessage(m string) Option {
	return func(e *Exception) { e.Message = m }
}

func WithError(err error, errs ...error) Option {
	if len(errs) > 0 {
		es := []error{err}
		es = append(es, errs...)
		err = errors.Join(es...)
	}
	return func(e *Exception) { e.Err = err }
}

func WithSeverity(s Severity) Option {
	return func(e *Exception) { e.Severity = s }
}

func WithData(key string, v any) Option {
	return func(e *Exception) {
		if e.Data == nil {
			e.Data = make(map[string]any)
		}
		e.Data[key] = v
	}
}

func WithHeader(header string, v string) Option {
	return func(e *Exception) {
		if e.headers == nil {
			e.headers = http.Header{}
		}
		e.headers.Add(header, v)
	}
}

func WithoutHeader(header string) Option {
	return func(e *Exception) {
		if e.headers == nil {
			e.headers = http.Header{}
		}
		e.headers.Del(header)
	}
}

func WithHandler(h HandlerFunc) Option {
	return func(e *Exception) { e.handler = h }
}
