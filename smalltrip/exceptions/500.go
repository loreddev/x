package exceptions

import (
	"errors"
	"net/http"
)

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
