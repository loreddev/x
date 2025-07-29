package problem

import (
	"encoding/xml"
	"errors"
	"net/http"
	"time"
)

func NewInternalError(err error, opts ...Option) InternalServerError {
	return InternalServerError{
		RegisteredProblem: NewDetailed(http.StatusInternalServerError, err.Error(), opts...),
		Errors:            newErrorTree(err).Errors,
		error:             err,
	}
}

type InternalServerError struct {
	RegisteredProblem
	Errors []ErrorTree `json:"errors" xml:"errors"`

	error error `json:"-" xml:"-"`
}

var _ error = InternalServerError{}

func (p InternalServerError) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Handler(p).ServeHTTP(w, r)
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
	return NotImplemented[T]{RegisteredProblem: p, RetryAfter: RetryAfter[T]{time: retryAfter}}
}

type NotImplemented[T time.Time | time.Duration] struct {
	RegisteredProblem
	RetryAfter RetryAfter[T] `json:"retryAfter" xml:"retry-after"`
}

func (p NotImplemented[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", p.RetryAfter.String())
	p.Handler(p).ServeHTTP(w, r)
}

func NewBadGateway(opts ...Option) BadGateway {
	return BadGateway{NewStatus(http.StatusBadGateway, opts...)}
}

type BadGateway struct{ RegisteredProblem }

func NewServiceUnavailable[T time.Time | time.Duration](retryAfter T, opts ...Option) ServiceUnavailable[T] {
	p := NewStatus(http.StatusNotImplemented, opts...)
	return ServiceUnavailable[T]{RegisteredProblem: p, RetryAfter: RetryAfter[T]{time: retryAfter}}
}

type ServiceUnavailable[T time.Time | time.Duration] struct {
	RegisteredProblem
	RetryAfter RetryAfter[T] `json:"retryAfter" xml:"retry-after"`
}

func (p ServiceUnavailable[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Retry-After", p.RetryAfter.String())
	p.handler(p).ServeHTTP(w, r)
}

func NewGatewayTimeout(opts ...Option) GatewayTimeout {
	return GatewayTimeout{NewStatus(http.StatusGatewayTimeout, opts...)}
}

type GatewayTimeout struct{ RegisteredProblem }

func NewHTTPVersionNotSupported(opts ...Option) HTTPVersionNotSupported {
	return HTTPVersionNotSupported{NewStatus(http.StatusHTTPVersionNotSupported, opts...)}
}

type HTTPVersionNotSupported struct{ RegisteredProblem }

func NewVariantAlsoNegotiates(opts ...Option) VariantAlsoNegotiates {
	return VariantAlsoNegotiates{NewStatus(http.StatusVariantAlsoNegotiates, opts...)}
}

type VariantAlsoNegotiates struct{ RegisteredProblem }

func NewInsufficientStorage(opts ...Option) InsufficientStorage {
	return InsufficientStorage{NewStatus(http.StatusInsufficientStorage, opts...)}
}

type InsufficientStorage struct{ RegisteredProblem }

func NewLoopDetected(opts ...Option) LoopDetected {
	return LoopDetected{NewStatus(http.StatusLoopDetected, opts...)}
}

type LoopDetected struct{ RegisteredProblem }

func NewNotExtended(opts ...Option) NotExtended {
	return NotExtended{NewStatus(http.StatusNotExtended, opts...)}
}

type NotExtended struct{ RegisteredProblem }

func NewNetworkAuthenticationRequired(opts ...Option) NetworkAuthenticationRequired {
	return NetworkAuthenticationRequired{NewStatus(http.StatusNetworkAuthenticationRequired, opts...)}
}

type NetworkAuthenticationRequired struct{ RegisteredProblem }
