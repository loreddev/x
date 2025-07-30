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

package problem

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func NewBadRequest(detail string, opts ...Option) BadRequest {
	return BadRequest{NewDetailed(http.StatusBadRequest, detail, opts...)}
}

type BadRequest struct{ RegisteredProblem }

func NewUnauthorized(scheme AuthScheme, opts ...Option) Unauthorized {
	return Unauthorized{
		RegisteredProblem: NewDetailed(
			http.StatusUnauthorized,
			fmt.Sprintf("You must authenticate using the %q scheme", scheme.Title),
			opts...,
		),
		Authentication: scheme,
	}
}

type Unauthorized struct {
	RegisteredProblem
	Authentication AuthScheme `json:"authentication,omitempty" xml:"authentication,omitempty"`
}

func (p Unauthorized) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Authentication.Title != "" {
		w.Header().Set("WWW-Authenticate", p.Authentication.Title)
	}
	p.Handler(p).ServeHTTP(w, r)
}

func NewPaymentRequired(opts ...Option) PaymentRequired {
	return PaymentRequired{NewStatus(http.StatusPaymentRequired, opts...)}
}

type PaymentRequired struct{ RegisteredProblem }

func NewForbidden(opts ...Option) Forbidden {
	return Forbidden{NewStatus(http.StatusForbidden, opts...)}
}

type Forbidden struct{ RegisteredProblem }

func NewNotFound(opts ...Option) NotFound {
	return NotFound{NewStatus(http.StatusNotFound, opts...)}
}

type NotFound struct{ RegisteredProblem }

func NewMethodNotAllowed[T string | []string](allow T, opts ...Option) MethodNotAllowed {
	p := MethodNotAllowed{
		RegisteredProblem: NewStatus(http.StatusMethodNotAllowed, opts...),
	}
	if as, ok := any(allow).([]string); ok {
		p.Allowed = as
	} else {
		p.Allowed = []string{any(allow).(string)}
	}
	return p
}

type MethodNotAllowed struct {
	RegisteredProblem
	Allowed []string `json:"allowed,omitempty" xml:"allowed,omitempty"`
}

func (p MethodNotAllowed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(p.Allowed) > 0 {
		w.Header().Set("Allow", strings.Join(p.Allowed, ", "))
	}
	p.Handler(p).ServeHTTP(w, r)
}

func NewNotAcceptable[T string | []string](header NegotiationHeader, allow T, opts ...Option) NotAcceptable {
	p := NotAcceptable{
		RegisteredProblem: NewDetailed(http.StatusMethodNotAllowed, fmt.Sprintf(
			"Cannot provide a response matching the list of acceptable values defined by %q header", header),
			opts...,
		),
	}
	if as, ok := any(allow).([]string); ok {
		p.Allowed = as
	} else {
		p.Allowed = []string{any(allow).(string)}
	}
	return p
}

type NotAcceptable struct {
	RegisteredProblem
	Allowed []string `json:"allowed,omitempty" xml:"allowed,omitempty"`
}

func (p NotAcceptable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Handler(p).ServeHTTP(w, r)
}

const (
	NegotiationHeaderAccept         NegotiationHeader = "Accept"
	NegotiationHeaderAcceptEncoding NegotiationHeader = "AcceptEncoding"
	NegotiationHeaderAcceptLanguage NegotiationHeader = "AcceptLanguage"
)

type NegotiationHeader string

func NewProxyAuthRequired(scheme AuthScheme, opts ...Option) ProxyAuthRequired {
	return ProxyAuthRequired{
		RegisteredProblem: NewDetailed(
			http.StatusProxyAuthRequired,
			fmt.Sprintf("You must authenticate on the proxy using the %q scheme", scheme.Title),
			opts...,
		),
		Authentication: scheme,
	}
}

type ProxyAuthRequired struct {
	RegisteredProblem
	Authentication AuthScheme `json:"authentication,omitempty" xml:"authentication,omitempty"`
}

func (p ProxyAuthRequired) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Authentication.Title != "" {
		w.Header().Set("Proxy-Authenticate", p.Authentication.Title)
	}
	p.Handler(p).ServeHTTP(w, r)
}

var (
	AuthSchemeBasic          = AuthScheme{Title: "Basic", Type: "https://datatracker.ietf.org/doc/html/rfc7617"}
	AuthSchemeBearer         = AuthScheme{Title: "Bearer", Type: "https://datatracker.ietf.org/doc/html/rfc6750"}
	AuthSchemeDigest         = AuthScheme{Title: "Digest", Type: "https://datatracker.ietf.org/doc/html/rfc7616"}
	AuthSchemeHOBA           = AuthScheme{Title: "HOBA", Type: "https://datatracker.ietf.org/doc/html/rfc7486"}
	AuthSchemeMutual         = AuthScheme{Title: "Mutual", Type: "https://datatracker.ietf.org/doc/html/rfc8120"}
	AuthSchemeNegotiate      = AuthScheme{Title: "Negotiate", Type: "https://datatracker.ietf.org/doc/html/rfc4599"}
	AuthSchemeVAPID          = AuthScheme{Title: "VAPID", Type: "https://datatracker.ietf.org/doc/html/rfc8292"}
	AuthSchemeSCRAM          = AuthScheme{Title: "SCRAM", Type: "https://datatracker.ietf.org/doc/html/rfc8292"}
	AuthSchemeAWS4HMACSHA256 = AuthScheme{Title: "AWS4-HMAC-SHA256", Type: "https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html"}
)

type AuthScheme struct {
	Type  string `json:"type,omitempty" xml:"type,omitempty"`
	Title string `json:"title,omitempty" xml:"title,omitempty"`
}

func NewRequestTimeout(opts ...Option) RequestTimeout {
	return RequestTimeout{NewStatus(http.StatusRequestTimeout, opts...)}
}

type RequestTimeout struct{ RegisteredProblem }

func (p RequestTimeout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "close")
	p.Handler(p).ServeHTTP(w, r)
}

func NewConflict(opts ...Option) Conflict {
	return Conflict{NewStatus(http.StatusConflict, opts...)}
}

type Conflict struct{ RegisteredProblem }

func NewGone(opts ...Option) Gone {
	return Gone{NewStatus(http.StatusGone, opts...)}
}

type Gone struct{ RegisteredProblem }

func NewLengthRequired(opts ...Option) LengthRequired {
	return LengthRequired{NewStatus(http.StatusLengthRequired, opts...)}
}

type LengthRequired struct{ RegisteredProblem }

func NewPreconditionFailed(opts ...Option) PreconditionFailed {
	return PreconditionFailed{NewStatus(http.StatusPreconditionFailed, opts...)}
}

type PreconditionFailed struct{ RegisteredProblem }

func NewContentTooLarge(opts ...Option) ContentTooLarge {
	return ContentTooLarge{NewStatus(http.StatusRequestEntityTooLarge, opts...)}
}

type ContentTooLarge struct{ RegisteredProblem }

func NewURITooLong(opts ...Option) URITooLong {
	return URITooLong{NewStatus(http.StatusRequestURITooLong, opts...)}
}

type URITooLong struct{ RegisteredProblem }

func NewUnsupportedMediaType(opts ...Option) UnsupportedMediaType {
	return UnsupportedMediaType{NewStatus(http.StatusUnsupportedMediaType, opts...)}
}

type UnsupportedMediaType struct{ RegisteredProblem }

func NewRangeNotSatisfiable(unit string, contentRange int, opts ...Option) RangeNotSatisfiable {
	return RangeNotSatisfiable{
		RegisteredProblem: NewStatus(http.StatusGone, opts...),
		Range:             fmt.Sprintf("%s */%d", unit, contentRange),
	}
}

type RangeNotSatisfiable struct {
	RegisteredProblem
	Range string `json:"range" xml:"range"`
}

func (p RangeNotSatisfiable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Range", p.Range)
	p.Handler(p).ServeHTTP(w, r)
}

func NewExpectationFailed(opts ...Option) ExpectationFailed {
	return ExpectationFailed{NewStatus(http.StatusExpectationFailed, opts...)}
}

type ExpectationFailed struct{ RegisteredProblem }

func NewTeapot(opts ...Option) Teapot {
	return Teapot{NewStatus(http.StatusTeapot, opts...)}
}

type Teapot struct{ RegisteredProblem }

func NewMisdirectedRequest(opts ...Option) MisdirectedRequest {
	return MisdirectedRequest{NewStatus(http.StatusMisdirectedRequest, opts...)}
}

type MisdirectedRequest struct{ RegisteredProblem }

func NewUnprocessableContent(opts ...Option) UnprocessableContent {
	return UnprocessableContent{NewStatus(http.StatusUnprocessableEntity, opts...)}
}

type UnprocessableContent struct{ RegisteredProblem }

// TODO?: Should the response of this be different and follow WebDAV's XML format?
func NewLocked(opts ...Option) Locked {
	return Locked{NewStatus(http.StatusLocked, opts...)}
}

type Locked struct{ RegisteredProblem }

func NewFailedDependency(opts ...Option) FailedDependency {
	return FailedDependency{NewStatus(http.StatusFailedDependency, opts...)}
}

type FailedDependency struct{ RegisteredProblem }

func NewTooEarly(opts ...Option) TooEarly {
	return TooEarly{NewStatus(http.StatusTooEarly, opts...)}
}

type TooEarly struct{ RegisteredProblem }

func NewUpgradeRequired(protocol Protocol, opts ...Option) UpgradeRequired {
	return UpgradeRequired{
		RegisteredProblem: NewStatus(http.StatusUpgradeRequired, opts...),
		Protocol:          protocol,
	}
}

type UpgradeRequired struct {
	RegisteredProblem
	Protocol Protocol `json:"protocol" xml:"protocol"`
}

func (p UpgradeRequired) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Upgrade", string(p.Protocol))
	p.Handler(p).ServeHTTP(w, r)
}

const (
	Http1_1 Protocol = "HTTP/1.1"
	Http2_0 Protocol = "HTTP/2.0"
	Http3_0 Protocol = "HTTP/3.0"
)

type Protocol string

func NewPreconditionRequired(opts ...Option) PreconditionRequired {
	return PreconditionRequired{NewStatus(http.StatusPreconditionRequired, opts...)}
}

type PreconditionRequired struct{ RegisteredProblem }

func NewTooManyRequests[T time.Time | time.Duration](retryAfter T, opts ...Option) TooManyRequests[T] {
	p := NewStatus(http.StatusTooManyRequests, opts...)
	return TooManyRequests[T]{RegisteredProblem: p, RetryAfter: RetryAfter[T]{time: retryAfter}}
}

type TooManyRequests[T time.Time | time.Duration] struct {
	RegisteredProblem
	RetryAfter RetryAfter[T] `json:"retryAfter" xml:"retry-after"`
}

func (p TooManyRequests[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", p.RetryAfter.String())
	p.Handler(p).ServeHTTP(w, r)
}

func NewRequestHeaderFieldsTooLarge(opts ...Option) RequestHeaderFieldsTooLarge {
	return RequestHeaderFieldsTooLarge{NewStatus(http.StatusRequestHeaderFieldsTooLarge, opts...)}
}

type RequestHeaderFieldsTooLarge struct{ RegisteredProblem }

func NewUnavailableForLegalReasons(opts ...Option) UnavailableForLegalReasons {
	return UnavailableForLegalReasons{NewStatus(http.StatusUnavailableForLegalReasons, opts...)}
}

type UnavailableForLegalReasons struct{ RegisteredProblem }
