package problem

import (
	"fmt"
	"net/http"
	"strings"
)

func NewBadRequest(detail string, opts ...Option) BadRequest {
	return NewDetailed(http.StatusBadRequest, detail, opts...)
}

type BadRequest = Problem

func NewUnauthorized(scheme AuthenticationScheme, opts ...Option) Unathorized {
	return Unathorized{
		Problem: NewDetailed(
			http.StatusUnauthorized,
			fmt.Sprintf("You must authenticate using the %q scheme", scheme.Title),
			opts...,
		),
		Authentication: scheme,
	}
}

type Unathorized struct {
	Problem
	Authentication AuthenticationScheme `json:"authentication,omitempty" xml:"authentication,omitempty"`
}

func (p Unathorized) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.Authentication.Title != "" {
		w.Header().Set("WWW-Authenticate", p.Authentication.Title)
	}
	p.handler(p).ServeHTTP(w, r)
}

var (
	AuthenticationSchemeBasic          = AuthenticationScheme{Title: "Basic", Type: "https://datatracker.ietf.org/doc/html/rfc7617"}
	AuthenticationSchemeBearer         = AuthenticationScheme{Title: "Bearer", Type: "https://datatracker.ietf.org/doc/html/rfc6750"}
	AuthenticationSchemeDigest         = AuthenticationScheme{Title: "Digest", Type: "https://datatracker.ietf.org/doc/html/rfc7616"}
	AuthenticationSchemeHOBA           = AuthenticationScheme{Title: "HOBA", Type: "https://datatracker.ietf.org/doc/html/rfc7486"}
	AuthenticationSchemeMutual         = AuthenticationScheme{Title: "Mutual", Type: "https://datatracker.ietf.org/doc/html/rfc8120"}
	AuthenticationSchemeNegotiate      = AuthenticationScheme{Title: "Negotiate", Type: "https://datatracker.ietf.org/doc/html/rfc4599"}
	AuthenticationSchemeVAPID          = AuthenticationScheme{Title: "VAPID", Type: "https://datatracker.ietf.org/doc/html/rfc8292"}
	AuthenticationSchemeSCRAM          = AuthenticationScheme{Title: "SCRAM", Type: "https://datatracker.ietf.org/doc/html/rfc8292"}
	AuthenticationSchemeAWS4HMACSHA256 = AuthenticationScheme{Title: "AWS4-HMAC-SHA256", Type: "https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html"}
)

type AuthenticationScheme struct {
	Type  string `json:"type,omitempty" xml:"type,omitempty"`
	Title string `json:"title,omitempty" xml:"title,omitempty"`
}

func NewPaymentRequired(opts ...Option) PaymentRequired {
	return NewStatus(http.StatusPaymentRequired, opts...)
}

type PaymentRequired = Problem

func NewForbidden(opts ...Option) Forbidden {
	return NewStatus(http.StatusForbidden, opts...)
}

type Forbidden = Problem

func NewNotFound(opts ...Option) NotFound {
	return NewStatus(http.StatusNotFound, opts...)
}

type NotFound = Problem

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
