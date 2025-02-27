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
	"net/http"
	"time"
)

// InternalServerError creates a new [Exception] with the "500 Internal Server Error"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
//
// An error should be provided to add context to the exception.
//
//	// "The HTTP 500 Internal Server Error server error response status
//	// code indicates that the server encountered an unexpected condition
//	// that prevented it from fulfilling the request. This error is a generic
//	// "catch-all" response to server issues, indicating that the server
//	// cannot find a more appropriate 5XX error to respond with."
//	//
//	// - Quoted from "500 Internal Server Error" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/502
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
//
//	// "The HTTP 501 Not Implemented server error response status
//	// code means that the server does not support the functionality
//	// required to fulfill the request.
//	//
//	// A response with this status may also include a Retry-After header,
//	// telling the client that they can retry the request after the specified
//	// time has elapsed. A 501 response is cacheable by default unless
//	// caching headers instruct otherwise.
//	//
//	// 501 is the appropriate response when the server does not recognize
//	// the request method and is incapable of supporting it for any resource.
//	// Servers are required to support GET and HEAD, and therefore must not
//	// return 501 in response to requests with these methods. If the server
//	// does recognize the method, but intentionally does not allow it, the
//	// appropriate response is 405 Method Not Allowed."
//	//
//	// - Quoted from "501 Not Implemented" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/502
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
//
//	// "The HTTP 502 Bad Gateway server error response status code
//	// indicates that a server was acting as a gateway or proxy and
//	// that it received an invalid response from the upstream server.
//	//
//	// This response is similar to a 500 Internal Server Error response
//	// in the sense that it is a generic "catch-call" for server errors.
//	// The difference is that it is specific to the point in the request
//	// chain that the error has occurred. If the origin server sends a
//	// valid HTTP error response to the gateway, the response should be
//	// passed on to the client instead of a 502 to make the failure
//	// reason transparent. If the proxy or gateway did not receive any
//	// HTTP response from the origin, it instead sends a
//	// 504 Gateway Timeout to the client."
//	//
//	// - Quoted from "502 Bad Gateway" by Mozilla Contributors, licensed
//	// under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/502
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
//
// A Retry-After header is passed with the duration provided by the "retryAfter"
// parameter.
//
//	// "The HTTP 503 Service Unavailable server error response status code
//	// indicates that the server is not ready to handle the request.
//	//
//	// Common causes are that a server is down for maintenance or overloaded.
//	// During maintenance, server administrators may temporarily route all traffic
//	// to a 503 page, or this may happen automatically during software updates.
//	// In overload cases, some server-side applications will reject requests with
//	// a 503 status when resource thresholds like memory, CPU, or connection pool
//	// limits are met. Dropping incoming requests creates backpressure that prevents
//	// the server's compute resources from being exhausted, avoiding more severe
//	// failures. If requests from specific clients are being restricted due to
//	// rate limiting, the appropriate response is 429 Too Many Requests.
//	//
//	// This response should be used for temporary conditions and the Retry-After
//	// HTTP header should contain the estimated time for the recovery of the
//	// service, if possible.
//	//
//	// A user-friendly page explaining the problem should be sent along with
//	// this response."
//	//
//	// - Quoted from "503 Service Unavailable" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/503
func ServiceUnavailable(retryAfter time.Time, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusServiceUnavailable),
		WithCode("Service Unavailable"),
		WithMessage("Not ready to handle the request."),
		WithError(errors.New("server is not ready to handle the request")),
		WithSeverity(ERROR),

		WithHeader("Retry-After", retryAfter.Format("Mon, 02 Jan 2006 15:04:05 GMT")),
	}
	o = append(o, opts...)

	return newException(o...)
}

// GatewayTimeout creates a new [Exception] with the "504 Gateway Timeout"
// status code, a human readable message and the provided error describing what in
// the request was wrong. The severity of this Exception by default is [ERROR].
//
//	// "The HTTP 504 Gateway Timeout server error response status code
//	// indicates that the server, while acting as a gateway or proxy,
//	// did not get a response in time from the upstream server in order
//	// to complete the request. This is similar to a 502 Bad Gateway,
//	// except that in a 504 status, the proxy or gateway did not receive
//	// any HTTP response from the origin within a certain time.
//	//
//	// There are many causes of 504 errors, and fixing such problems
//	// likely requires investigation and debugging by server administrators,
//	// or the site may work again at a later time. Exceptions are client
//	// networking errors, particularly if the service works for other
//	// visitors, and if clients use VPNs or other custom networking setups.
//	// In such cases, clients should check network settings, firewall setup,
//	// proxy settings, DNS configuration, etc."
//	//
//	// - Quoted from "504 Gateway Timeout" by Mozilla Contributors, licensed
//	// under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/504
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
//
//	// "The HTTP 505 HTTP Version Not Supported server error response status
//	// code indicates that the HTTP version used in the request is not supported
//	// by the server.
//	//
//	// It's common to see this error when a request line is improperly formed
//	// such as GET /path to resource HTTP/1.1 or with \n terminating the request
//	// line instead of \r\n."
//	//
//	// - Quoted from "505 HTTP Version Not Supported" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/505
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
//
//	// "The HTTP 506 Variant Also Negotiates server error response status code
//	// is returned during content negotiation when there is recursive loop in
//	// the process of selecting a resource.
//	//
//	// Agent-driven content negotiation enables a client and server to
//	// collaboratively decide the best variant of a given resource when the server
//	// has multiple variants. A server sends a 506 status code due to server
//	// misconfiguration that results in circular references when creating responses."
//	//
//	// - Quoted from "506 Variant Also Negotiates" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/506
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
//
//	// "The HTTP 507 Insufficient Storage server error response status code
//	// indicates that an action could not be performed because the server does
//	// not have enough available storage to successfully complete the request.
//	//
//	// [...] Common causes of this error can be from server directories running
//	// out of available space, not enough available RAM for an operation, or
//	// internal limits reached (such as application-specific memory limits,
//	// for example). The request causing this error does not necessarily need to
//	// include content, as it may be a request that would create a resource on
//	// the server if it was successful."
//	//
//	// - Quoted from "507 Insufficient Storage" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/507
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
//
//	// "The HTTP 508 Loop Detected server error response status code indicates
//	// that the entire operation failed because it encountered an infinite loop
//	// while processing a request with Depth: infinity.
//	//
//	// The status may be given in the context of the Web Distributed Authoring
//	// and Versioning (WebDAV). It was introduced as a fallback for cases where
//	// WebDAV clients do not support 208 Already Reported responses (when requests
//	// do not explicitly include a DAV header)."
//	//
//	// - Quoted from "508 Loop Detected" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/508
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
//
//	// "The HTTP 510 Not Extended server error response status code is sent
//	// when the client request declares an HTTP Extension (RFC 2774) that
//	// should be used to process the request, but the extension is not supported."
//	//
//	// - Quoted from "510 Not Extended" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/510
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
//
//	// "The HTTP 511 Network Authentication Required server error response status
//	// code indicates that the client needs to authenticate to gain network access.
//	// This status is not generated by origin servers, but by intercepting proxies
//	// that control access to a network.
//	//
//	// Network operators sometimes require some authentication, acceptance of terms,
//	// or other user interaction before granting access (for example in an internet
//	// café or at an airport). They often identify clients who have not done so using
//	// their Media Access Control (MAC) addresses."
//	//
//	// - Quoted from "511 Network Authentication Required" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/511
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
