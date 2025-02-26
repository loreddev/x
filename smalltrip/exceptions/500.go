package exceptions

import (
	"errors"
	"net/http"
)

// InternalServerError creates a new [Exception] with the "500 Internal Server Error"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
//
// An error should be provided to add context to the exception.
func InternalServerError(err error, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusInternalServerError),
		WithCode("Internal Server Error"),
		WithMessage("A unexpected error occurred."),
		WithError(err),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// NotImplemented creates a new [Exception] with the "501 Not Implemented"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func NotImplemented(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusNotImplemented),
		WithCode("Not Implemented"),
		WithMessage("Functionality is not supported."),
		WithError(errors.New("user agent requested functionality that is not supported")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// BadGateway creates a new [Exception] with the "502 Bad Gateway"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func BadGateway(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusBadGateway),
		WithCode("Bad Gateway"),
		WithMessage("Invalid response from upstream."),
		WithError(errors.New("upstream response is invalid")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// ServiceUnavailable creates a new [Exception] with the "503 Service Unavailable"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func ServiceUnavailable(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusServiceUnavailable),
		WithCode("Service Unavailable"),
		WithMessage("Not ready to handle the request."),
		WithError(errors.New("server is not ready to handle the request")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// GatewayTimeout creates a new [Exception] with the "504 Gateway Timeout"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func GatewayTimeout(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusGatewayTimeout),
		WithCode("Gateway Timeout"),
		WithMessage("Did not get the response from upstream in time."),
		WithError(errors.New("upstream did not respond in time")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// HTTPVersionNotSupported creates a new [Exception] with the "505 HTTP Version Not Supported"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func HTTPVersionNotSupported(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusHTTPVersionNotSupported),
		WithCode("HTTP Version Not Supported"),
		WithMessage("Version of HTTP is not supported."),
		WithError(errors.New("server does not support requested HTTP version")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// VariantAlsoNegotiates creates a new [Exception] with the "506 Variant Also Negotiates"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func VariantAlsoNegotiates(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusVariantAlsoNegotiates),
		WithCode("Variant Also Negotiates"),
		WithMessage("A recursive loop found in the process of the request."),
		WithError(errors.New("variant also negotiates, recursive loop found on request")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// InsufficientStorage creates a new [Exception] with the "507 Insufficient Storage"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func InsufficientStorage(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusInsufficientStorage),
		WithCode("Insufficient Storage"),
		WithMessage("There is not enough available storage to complete the request."),
		WithError(errors.New("not enough available storage to complete request")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// LoopDetected creates a new [Exception] with the "508 Loop Detected"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func LoopDetected(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusLoopDetected),
		WithCode("Loop Detected"),
		WithMessage("Infinite loop found while processing the request."),
		WithError(errors.New("infinite loop found while processing request")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// NotExtended creates a new [Exception] with the "510 Not Extended"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func NotExtended(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusLoopDetected),
		WithCode("Not Extended"),
		WithMessage("HTTP extension is not supported."),
		WithError(errors.New("user agent requested with a HTTP extension that is not supported")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}

// NetworkAuthenticationRequired creates a new [Exception] with the "511 Network Authentication Required"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
func NetworkAuthenticationRequired(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusNetworkAuthenticationRequired),
		WithCode("Network Authentication Required"),
		WithMessage("Authentication to access network access is necessary."),
		WithError(errors.New("user agent requested without being network authenticated")),
		WithSeverity(ERROR),
	}
	o = append(o, opts...)

	return newException(o...)
}
