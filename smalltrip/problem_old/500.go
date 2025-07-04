package problem

import (
	"fmt"
	"net/http"
	"time"

	"forge.capytal.company/loreddev/x/smalltrip/problem/extension"
)

func InternalServerError(err error, opts ...any) Problem {
	p := internalServerError{
		RegisteredMembers: RegisteredMembers{
			TypeURI:       "about:blank",
			TypeTitle:     "Internal Server Error",
			StatusCode:    http.StatusInternalServerError,
			DetailMessage: err.Error(),
		},
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case RegisteredMembersOption:
			opt(&p.RegisteredMembers)
		case Option:
			opt(&p)
		}
	}

	p.Errors = extension.NewErrorTree(err)
	return p
}

type internalServerError struct {
	RegisteredMembers
	extension.Errors
}

func (p internalServerError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeProblem(p, w, r)
}

func NotImplemented(opts ...any) Problem {
	p := notImplemented{
		RegisteredMembers: RegisteredMembers{
			TypeURI:    "about:blank",
			TypeTitle:  "Not Implemented",
			StatusCode: http.StatusNotImplemented,
		},
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case RegisteredMembersOption:
			opt(&p.RegisteredMembers)
		case Option:
			opt(&p)
		}
	}

	return p
}

type notImplemented struct {
	RegisteredMembers
	retryAfterDuration time.Duration `json:"-" xml:"-"`
	retryAfterTime     time.Time     `json:"-" xml:"-"`
}

var _ hasRetryAfter = (*notImplemented)(nil)

func (p *notImplemented) SetRetryAfterTime(t time.Time) {
	p.retryAfterTime = t
}

func (p *notImplemented) SetRetryAfterDuration(d time.Duration) {
	p.retryAfterDuration = d
}

func (p notImplemented) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !p.retryAfterTime.IsZero() {
		w.Header().Set("Retry-After", p.retryAfterTime.UTC().Format(http.TimeFormat))
	} else if p.retryAfterDuration != 0 {
		w.Header().Set("Retry-After", fmt.Sprintf("%.0f", p.retryAfterDuration.Seconds()))
	}

	ServeProblem(p, w, r)
}

func BadGateway(opts ...any) Problem {
	p := badGateway{
		RegisteredMembers: RegisteredMembers{
			TypeURI:    "about:blank",
			TypeTitle:  "Bad Gateway",
			StatusCode: http.StatusBadGateway,
		},
	}

	for _, opt := range opts {
		switch opt := opt.(type) {
		case RegisteredMembersOption:
			opt(&p.RegisteredMembers)
		case Option:
			opt(&p)
		}
	}

	return p
}

type badGateway struct{ RegisteredMembers }

func (p badGateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeProblem(p, w, r)
}
