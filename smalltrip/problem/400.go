package problem

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func NewBadRequest(detail string, opts ...Option) BadRequest {
	return BadRequest(NewDetailed(http.StatusBadRequest, detail, opts...))
}

type BadRequest Problem

func NewUnauthorized(scheme AuthScheme, opts ...Option) Unauthorized {
	return Unauthorized{
		Problem: NewDetailed(
			http.StatusUnauthorized,
			fmt.Sprintf("You must authenticate using the %q scheme", scheme.Title),
			opts...,
		),
		Authentication: scheme,
	}
}

type Unauthorized struct {
	Problem
	Authentication AuthScheme `json:"authentication,omitempty" xml:"authentication,omitempty"`
}

func (p Unauthorized) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Authentication.Title != "" {
		w.Header().Set("WWW-Authenticate", p.Authentication.Title)
	}
	p.handler(p).ServeHTTP(w, r)
}

func NewPaymentRequired(opts ...Option) PaymentRequired {
	return PaymentRequired(NewStatus(http.StatusPaymentRequired, opts...))
}

type PaymentRequired Problem

func NewForbidden(opts ...Option) Forbidden {
	return Forbidden(NewStatus(http.StatusForbidden, opts...))
}

type Forbidden Problem

func NewNotFound(opts ...Option) NotFound {
	return NotFound(NewStatus(http.StatusNotFound, opts...))
}

type NotFound Problem

func NewMethodNotAllowed[T string | []string](allow T, opts ...Option) MethodNotAllowed {
	p := MethodNotAllowed{
		Problem: NewStatus(http.StatusMethodNotAllowed, opts...),
	}
	if as, ok := any(allow).([]string); ok {
		p.Allowed = as
	} else {
		p.Allowed = []string{any(allow).(string)}
	}
	return p
}

type MethodNotAllowed struct {
	Problem
	Allowed []string `json:"allowed,omitempty" xml:"allowed,omitempty"`
}

func (p MethodNotAllowed) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(p.Allowed) > 0 {
		w.Header().Set("Allow", strings.Join(p.Allowed, ", "))
	}
	p.handler(p).ServeHTTP(w, r)
}

func NewNotAcceptable[T string | []string](header NegotiationHeader, allow T, opts ...Option) NotAcceptable {
	p := NotAcceptable{
		Problem: NewDetailed(http.StatusMethodNotAllowed, fmt.Sprintf(
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
	Problem
	Allowed []string `json:"allowed,omitempty" xml:"allowed,omitempty"`
}

func (p NotAcceptable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.handler(p).ServeHTTP(w, r)
}

const (
	NegotiationHeaderAccept         NegotiationHeader = "Accept"
	NegotiationHeaderAcceptEncoding NegotiationHeader = "AcceptEncoding"
	NegotiationHeaderAcceptLanguage NegotiationHeader = "AcceptLanguage"
)

type NegotiationHeader string

func NewProxyAuthRequired(scheme AuthScheme, opts ...Option) ProxyAuthRequired {
	return ProxyAuthRequired{
		Problem: NewDetailed(
			http.StatusProxyAuthRequired,
			fmt.Sprintf("You must authenticate on the proxy using the %q scheme", scheme.Title),
			opts...,
		),
		Authentication: scheme,
	}
}

type ProxyAuthRequired struct {
	Problem
	Authentication AuthScheme `json:"authentication,omitempty" xml:"authentication,omitempty"`
}

func (p ProxyAuthRequired) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Authentication.Title != "" {
		w.Header().Set("Proxy-Authenticate", p.Authentication.Title)
	}
	p.handler(p).ServeHTTP(w, r)
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
	return RequestTimeout(NewStatus(http.StatusRequestTimeout, opts...))
}

type RequestTimeout Problem

func (p RequestTimeout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Connection", "close")
	p.handler(p).ServeHTTP(w, r)
}

func NewConflict(opts ...Option) Conflict {
	return Conflict(NewStatus(http.StatusConflict, opts...))
}

type Conflict Problem

func NewGone(opts ...Option) Gone {
	return Gone(NewStatus(http.StatusGone, opts...))
}

type Gone Problem

func NewLengthRequired(opts ...Option) LengthRequired {
	return LengthRequired(NewStatus(http.StatusLengthRequired, opts...))
}

type LengthRequired Problem

func NewPreconditionFailed(opts ...Option) PreconditionFailed {
	return PreconditionFailed(NewStatus(http.StatusPreconditionFailed, opts...))
}

type PreconditionFailed Problem

func NewContentTooLarge(opts ...Option) ContentTooLarge {
	return ContentTooLarge(NewStatus(http.StatusRequestEntityTooLarge, opts...))
}

type ContentTooLarge Problem

func NewURITooLong(opts ...Option) URITooLong {
	return URITooLong(NewStatus(http.StatusRequestURITooLong, opts...))
}

type URITooLong Problem

func NewUnsupportedMediaType(opts ...Option) UnsupportedMediaType {
	return UnsupportedMediaType(NewStatus(http.StatusUnsupportedMediaType, opts...))
}

type UnsupportedMediaType Problem

func NewRangeNotSatisfiable(unit string, contentRange int, opts ...Option) RangeNotSatisfiable {
	return RangeNotSatisfiable{
		Problem: NewStatus(http.StatusGone, opts...),
		Range:   fmt.Sprintf("%s */%d", unit, contentRange),
	}
}

type RangeNotSatisfiable struct {
	Problem
	Range string `json:"range" xml:"range"`
}

func (p RangeNotSatisfiable) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Range", p.Range)
	p.handler(p).ServeHTTP(w, r)
}

func NewExpectationFailed(opts ...Option) ExpectationFailed {
	return ExpectationFailed(NewStatus(http.StatusExpectationFailed, opts...))
}

type ExpectationFailed Problem

func NewTeapot(opts ...Option) Teapot {
	return Teapot(NewStatus(http.StatusTeapot, opts...))
}

type Teapot Problem

func NewMisdirectedRequest(opts ...Option) MisdirectedRequest {
	return MisdirectedRequest(NewStatus(http.StatusMisdirectedRequest, opts...))
}

type MisdirectedRequest Problem

func NewUnprocessableContent(opts ...Option) UnprocessableContent {
	return UnprocessableContent(NewStatus(http.StatusUnprocessableEntity, opts...))
}

type UnprocessableContent Problem

// TODO?: Should the response of this be different and follow WebDAV's XML format?
func NewLocked(opts ...Option) Locked {
	return Locked(NewStatus(http.StatusLocked, opts...))
}

type Locked Problem

func NewFailedDependency(opts ...Option) FailedDependency {
	return FailedDependency(NewStatus(http.StatusFailedDependency, opts...))
}

type FailedDependency Problem

func NewTooEarly(opts ...Option) TooEarly {
	return TooEarly(NewStatus(http.StatusTooEarly, opts...))
}

type TooEarly Problem

func NewUpgradeRequired(protocol Protocol, opts ...Option) UpgradeRequired {
	return UpgradeRequired{
		Problem:  NewStatus(http.StatusUpgradeRequired, opts...),
		Protocol: protocol,
	}
}

type UpgradeRequired struct {
	Problem
	Protocol Protocol `json:"protocol" xml:"protocol"`
}

func (p UpgradeRequired) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Upgrade", string(p.Protocol))
	p.handler(p).ServeHTTP(w, r)
}

const (
	Http1_1 Protocol = "HTTP/1.1"
	Http2_0 Protocol = "HTTP/2.0"
	Http3_0 Protocol = "HTTP/3.0"
)

type Protocol string

func NewPreconditionRequired(opts ...Option) PreconditionRequired {
	return PreconditionRequired(NewStatus(http.StatusPreconditionRequired, opts...))
}

type PreconditionRequired Problem

func NewTooManyRequests[T time.Time | time.Duration](retryAfter T, opts ...Option) TooManyRequests[T] {
	p := NewStatus(http.StatusTooManyRequests, opts...)
	return TooManyRequests[T]{Problem: p, RetryAfter: retryAfter}
}

type TooManyRequests[T time.Time | time.Duration] struct {
	Problem
	RetryAfter RetryAfter[T] `json:"retryAfter" xml:"retry-after"`
}

func (p TooManyRequests[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", p.RetryAfter.String())
	p.handler(p).ServeHTTP(w, r)
}

func NewRequestHeaderFieldsTooLarge(opts ...Option) RequestHeaderFieldsTooLarge {
	return RequestHeaderFieldsTooLarge(NewStatus(http.StatusRequestHeaderFieldsTooLarge, opts...))
}

type RequestHeaderFieldsTooLarge Problem

func NewUnavailableForLegalReasons(opts ...Option) UnavailableForLegalReasons {
	return UnavailableForLegalReasons(NewStatus(http.StatusUnavailableForLegalReasons, opts...))
}

type UnavailableForLegalReasons Problem
