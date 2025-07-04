package problem

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func NewInternalError(err error, opts ...Option) InternalServerError {
	return InternalServerError{
		Problem: NewDetailed(http.StatusInternalServerError, err.Error(), opts...),
		Errors:  newErrorTree(err).Errors,
		error:   err,
	}
}

type InternalServerError struct {
	Problem
	Errors []ErrorTree `json:"errors" xml:"errors"`

	error error `json:"-" xml:"-"`
}

var _ error = InternalServerError{}

func (p InternalServerError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.handler(p).ServeHTTP(w, r)
}

func (p InternalServerError) Error() string {
	return p.error.Error()
}

func newErrorTree(err error) ErrorTree {
	i := ErrorTree{Detail: err.Error(), Errors: []ErrorTree{}, error: err}
	if us, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range us.Unwrap() {
			i.Errors = append(i.Errors, newErrorTree(e))
		}
	} else if e := errors.Unwrap(err); e != nil {
		i.Errors = append(i.Errors, newErrorTree(e))
	}
	return i
}

type ErrorTree struct {
	Detail string      `json:"detail" xml:"detail"`
	Errors []ErrorTree `json:"errors" xml:"errors"`

	XMLName xml.Name `json:"-" xml:"errors"`
	error   error    `json:"-" xml:"-"`
}

var _ error = ErrorTree{}

func (i ErrorTree) Error() string {
	return i.error.Error()
}

func NewNotImplemented[T time.Time | time.Duration](retryAfter T, opts ...Option) NotImplemented[T] {
	p := NewStatus(http.StatusNotImplemented, opts...)
	return NotImplemented[T]{Problem: p, RetryAfter: retryAfter}
}

type NotImplemented[T time.Time | time.Duration] struct {
	Problem
	RetryAfter T `json:"retryAfter" xml:"retry-after"`
}

var (
	_ json.Marshaler = NotImplemented[time.Time]{}
	_ xml.Marshaler  = NotImplemented[time.Time]{}
)

func (p NotImplemented[T]) MarshalJSON() ([]byte, error) {
	switch t := any(p.RetryAfter).(type) {
	case time.Time:
		return json.Marshal(struct {
			Problem
			RetryAfter string `json:"retryAfter,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: t.Format(time.RFC3339),
		})
	case time.Duration:
		return json.Marshal(struct {
			Problem
			RetryAfter int `json:"retryAfter,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: int(t.Seconds()),
		})
	default:
		return nil, errors.New("problems-not-implemented: RetryAfter is not of type time.Time or time.Duration")
	}
}

func (p NotImplemented[T]) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	switch t := any(p.RetryAfter).(type) {
	case time.Time:
		return e.Encode(struct {
			Problem
			RetryAfter string `xml:"retry-after,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: t.Format(time.RFC3339),
		})
	case time.Duration:
		return e.Encode(struct {
			Problem
			RetryAfter int `xml:"retry-after,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: int(t.Seconds()),
		})
	default:
		return errors.New("problems-not-implemented: RetryAfter is not of type time.Time or time.Duration")
	}
}

func (p NotImplemented[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch t := any(p.RetryAfter).(type) {
	case time.Time:
		if !t.IsZero() {
			w.Header().Set("Retry-After", t.Format(http.TimeFormat))
		}
	case time.Duration:
		if t != 0 {
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", t.Seconds()))
		}
	}
	p.handler(p).ServeHTTP(w, r)
}

func NewBadGateway(opts ...Option) BadGateway {
	return NewStatus(http.StatusBadGateway, opts...)
}

type BadGateway = Problem

func NewServiceUnavailable[T time.Time | time.Duration](retryAfter T, opts ...Option) ServiceUnavailable[T] {
	p := NewStatus(http.StatusNotImplemented, opts...)
	return ServiceUnavailable[T]{Problem: p, RetryAfter: retryAfter}
}

type ServiceUnavailable[T time.Time | time.Duration] struct {
	Problem
	RetryAfter T `json:"retryAfter" xml:"retry-after"`
}

var (
	_ json.Marshaler = ServiceUnavailable[time.Time]{}
	_ xml.Marshaler  = ServiceUnavailable[time.Time]{}
)

func (p ServiceUnavailable[T]) MarshalJSON() ([]byte, error) {
	switch t := any(p.RetryAfter).(type) {
	case time.Time:
		return json.Marshal(struct {
			Problem
			RetryAfter string `json:"retryAfter,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: t.Format(time.RFC3339),
		})
	case time.Duration:
		return json.Marshal(struct {
			Problem
			RetryAfter int `json:"retryAfter,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: int(t.Seconds()),
		})
	default:
		return nil, errors.New("problems-not-implemented: RetryAfter is not of type time.Time or time.Duration")
	}
}

func (p ServiceUnavailable[T]) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	switch t := any(p.RetryAfter).(type) {
	case time.Time:
		return e.Encode(struct {
			Problem
			RetryAfter string `xml:"retry-after,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: t.Format(time.RFC3339),
		})
	case time.Duration:
		return e.Encode(struct {
			Problem
			RetryAfter int `xml:"retry-after,omitempty"`
		}{
			Problem:    p.Problem,
			RetryAfter: int(t.Seconds()),
		})
	default:
		return errors.New("problems-not-implemented: RetryAfter is not of type time.Time or time.Duration")
	}
}

func (p ServiceUnavailable[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch t := any(p.RetryAfter).(type) {
	case time.Time:
		if !t.IsZero() {
			w.Header().Set("Retry-After", t.Format(http.TimeFormat))
		}
	case time.Duration:
		if t != 0 {
			w.Header().Set("Retry-After", fmt.Sprintf("%.0f", t.Seconds()))
		}
	}
	p.handler(p).ServeHTTP(w, r)
}

func NewGatewayTimeout(opts ...Option) GatewayTimeout {
	return NewStatus(http.StatusGatewayTimeout, opts...)
}

type GatewayTimeout = Problem

func NewHTTPVersionNotSupported(opts ...Option) HTTPVersionNotSupported {
	return NewStatus(http.StatusHTTPVersionNotSupported, opts...)
}

type HTTPVersionNotSupported = Problem

func NewVariantAlsoNegotiates(opts ...Option) VariantAlsoNegotiates {
	return NewStatus(http.StatusVariantAlsoNegotiates, opts...)
}

type VariantAlsoNegotiates = Problem

func NewInsufficientStorage(opts ...Option) InsufficientStorage {
	return NewStatus(http.StatusInsufficientStorage, opts...)
}

type InsufficientStorage = Problem

func NewLoopDetected(opts ...Option) LoopDetected {
	return NewStatus(http.StatusLoopDetected, opts...)
}

type LoopDetected = Problem

func NewNotExtended(opts ...Option) NotExtended {
	return NewStatus(http.StatusNotExtended, opts...)
}

type NotExtended = Problem

func NewNetworkAuthenticationRequired(opts ...Option) NetworkAuthenticationRequired {
	return NewStatus(http.StatusNetworkAuthenticationRequired, opts...)
}

type NetworkAuthenticationRequired = Problem
