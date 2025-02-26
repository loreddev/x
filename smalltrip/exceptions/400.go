package exceptions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// BadRequest creates a new [Exception] with the "400 Bad Request" status code,
// a human readable message and the provided error describing what in the request
// was wrong. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 400 Bad Request client error response status code indicates that the
//	// server would not process the request due to something the server considered to
//	// be a client error. The reason for a 400 response is typically due to malformed
//	// request syntax, invalid request message framing, or deceptive request routing.
//	//
//	// Clients that receive a 400 response should expect that repeating the request
//	// without modification will fail with the same error."
//	//
//	// - Quoted from "400 Bad Request" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/400
func BadRequest(err error, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusBadRequest),
		WithCode("Bad Request"),
		WithMessage("The request sent is malformed, see the provided error for more information."),
		WithError(err),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// Unathorized creates a new [Exception] with the "401 Unathorized" status code,
// a human readable message and error. The severity of this Exception by default
// is [WARN]. A "WWW-Authenticate" header should be sent with this exception,
// provided via the "authenticate" parameter.
//
//	// "The HTTP 401 Unauthorized client error response status code indicates that a
//	// request was not successful because it lacks valid authentication credentials
//	// for the requested resource. This status code is sent with an HTTP WWW-Authenticate
//	// response header that contains information on the authentication scheme the
//	// server expects the client to include to make the request successfully.
//	//
//	// A 401 Unauthorized is similar to the 403 Forbidden response, except that a 403
//	// is returned when a request contains valid credentials, but the client does not
//	// have permissions to perform a certain action."
//	//
//	// - Quoted from "401 Unathorized" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/401
func Unathorized(authenticate string, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusUnauthorized),
		WithCode("Unathorized"),
		WithMessage("Not authorized/authenticated to be able to do this request."),
		WithError(errors.New("user agent is not authenticated to do the request")),
		WithSeverity(WARN),

		WithHeader("WWW-Authenticate", authenticate),
	}
	o = append(o, opts...)

	return newException(o...)
}

// PaymentRequired creates a new [Exception] with the "402 Payment Required" status code,
// a human readable message and error. The severity of this Exception by default
// is [WARN].
//
//	// "The HTTP 402 Payment Required client error response status code is a nonstandard
//	// response status code reserved for future use.
//	//
//	// This status code was created to enable digital cash or (micro) payment systems
//	// and would indicate that requested content is not available until the client makes
//	// a payment. No standard use convention exists and different systems use it in
//	// different contexts.
//	//
//	// [...]
//	//
//	// This status code is reserved but not defined. Actual implementations vary in
//	// the format and contents of the response. No browser supports a 402, and an error
//	// will be displayed as a generic 4xx status code."
//	//
//	// - Quoted from "402 Payment Required" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/402
func PaymentRequired(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusPaymentRequired),
		WithCode("Payment Required"),
		WithMessage("Payment is required to be able to see this page."),
		WithError(errors.New("user agent needs to have payed to see this content")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// Forbidden creates a new [Exception] with the "403 Forbidden" status code,
// a human readable message and error. The severity of this Exception by default
// is [WARN].
//
//	// "The HTTP 403 Forbidden client error response status code indicates that the
//	// server understood the request but refused to process it. This status is similar
//	// to 401, except that for 403 Forbidden responses, authenticating or
//	// re-authenticating makes no difference. The request failure is tied to application
//	// logic, such as insufficient permissions to a resource or action.
//	//
//	// Clients that receive a 403 response should expect that repeating the request
//	// without modification will fail with the same error. Server owners may decide
//	// to send a 404 response instead of a 403 if acknowledging the existence of a
//	// resource to clients with insufficient privileges is not desired."
//	//
//	// - Quoted from "403 Forbidden" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/403
func Forbidden(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusForbidden),
		WithCode("Forbidden"),
		WithMessage("You do not have the rights to do this request."),
		WithError(errors.New("user agent does not have the rights to do the request")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// NotFound creates a new [Exception] with the "404 Not Found" status code,
// a human readable message and error. The severity of this Exception by default
// is [WARN].
//
//	// "The HTTP 404 Not Found client error response status code indicates that the
//	// server cannot find the requested resource. Links that lead to a 404 page are
//	// often called broken or dead links and can be subject to link rot.
//	//
//	// A 404 status code only indicates that the resource is missing without indicating
//	// if this is temporary or permanent. If a resource is permanently removed, servers
//	// should send the 410 Gone status instead.
//	//
//	// 404 errors on a website can lead to a poor user experience for your visitors,
//	// so the number of broken links (internal and external) should be minimized to
//	// prevent frustration for readers. Common causes of 404 responses are mistyped
//	// URLs or pages that are moved or deleted without redirection. For more
//	// information, see the Redirections in HTTP guide."
//	//
//	// - Quoted from "404 Not Found" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/404
func NotFound(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusNotFound),
		WithCode("Not Found"),
		WithMessage("This content does not exists."),
		WithError(errors.New("user agent requested content that does not exists")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// MethodNotAllowed creates a new [Exception] with the "405 Method Not Allowed" status code,
// a human readable message and error, and a "Allow" header with the methods provided via the
// "allowed" parameter. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 405 Method Not Allowed client error response status code indicates
//	// that the server knows the request method, but the target resource doesn't support
//	// this method. The server must generate an Allow header in a 405 response with a
//	// list of methods that the target resource currently supports.
//	//
//	// Improper server-side permissions set on files or directories may cause a 405
//	// response when the request would otherwise be expected to succeed."
//	//
//	// - Quoted from "405 Method Not Allowed" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/405
func MethodNotAllowed(allowed []string, opts ...Option) Exception {
	a := strings.Join(allowed, ", ")
	o := []Option{
		WithStatus(http.StatusMethodNotAllowed),
		WithCode("Method Not Allowed"),
		WithMessage("The method is not allowed for this endpoints."),
		WithError(fmt.Errorf("user agent tried to use method which is not a allowed method (%s)", a)),
		WithSeverity(WARN),

		WithHeader("Allow", a),
	}
	o = append(o, opts...)

	return newException(o...)
}

// NotAcceptable creates a new [Exception] with the "406 Not Acceptable" status code,
// a human readable message and error, and a list of accepted mime types provided via
// the "accepted" parameter. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 406 Not Acceptable client error response status code indicates that the
//	// server could not produce a response matching the list of acceptable values defined
//	// in the request's proactive content negotiation headers and that the server was
//	// unwilling to supply a default representation.
//	//
//	// [...]
//	//
//	// A server may return responses that differ from the request's accept headers.
//	// In such cases, a 200 response with a default resource that doesn't match the client's
//	// list of acceptable content negotiation values may be preferable to sending a 406 response.
//	//
//	// If a server returns a 406, the body of the message should contain the list of
//	// available representations for the resource, allowing the user to choose, although
//	// no standard way for this is defined."
//	//
//	// - Quoted from "406 Not Acceptable" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/406
func NotAcceptable(accepted []string, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusNotAcceptable),
		WithCode("Not Acceptable"),
		WithMessage("Unable to find any content that conforms to the request."),
		WithError(errors.New("no content conforms to the requested criteria of the user agent")),
		WithData("accepted", accepted),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// ProxyAuthenticationRequired creates a new [Exception] with the "407 Proxy Authentication Required"
// status code, a human readable message and error. The severity of this Exception by default is [WARN].
// A "Proxy-Authenticate" header should be sent with this exception, provided via the
// "authenticate" parameter.
//
//	// "The HTTP 407 Proxy Authentication Required client error response status code indicates that
//	// the request did not succeed because it lacks valid authentication credentials for the proxy
//	// server that sits between the client and the server with access to the requested resource.
//	//
//	// This response is sent with a Proxy-Authenticate header that contains information on how to
//	// correctly authenticate requests. The client may repeat the request with a new or replaced
//	// Proxy-Authorization header field."
//	//
//	// - Quoted from "407 Proxy Authentication Required" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/407
func ProxyAuthenticationRequired(authenticate string, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusProxyAuthRequired),
		WithCode("Proxy Authentication Required"),
		WithMessage("Authorization/authentication via proxy is needed to access this content."),
		WithError(errors.New("user agent is missing proxy authorization/authentication")),
		WithSeverity(WARN),

		WithHeader("Proxy-Authenticate", authenticate),
	}
	o = append(o, opts...)

	return newException(o...)
}

// RequestTimeout creates a new [Exception] with the "408 Request Timeout" status code, a human
// readable message and error, with a "Connection: close" header alongside. The severity of this
// Exception by default is [WARN].
//
//	// "The HTTP 408 Request Timeout client error response status code indicates that the server
//	// would like to shut down this unused connection. A 408 is sent on an idle connection by some
//	// servers, even without any previous request by the client.
//	//
//	// A server should send the Connection: close header field in the response, since 408 implies
//	// that the server has decided to close the connection rather than continue waiting.
//	//
//	// This response is used much more since some browsers, like Chrome and Firefox, use HTTP
//	// pre-connection mechanisms to speed up surfing."
//	//
//	// - Quoted from "408 Request Timeout" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/408
func RequestTimeout(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusRequestTimeout),
		WithCode("Request Timeout"),
		WithMessage("This request was shut down."),
		WithError(errors.New("request timed out")),
		WithSeverity(WARN),

		WithHeader("Connection", "close"),
	}
	o = append(o, opts...)

	return newException(o...)
}

// Conflict creates a new [Exception] with the "409 Conflict" status code, a human
// readable message and error. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 409 Conflict client error response status code indicates a request
//	// conflict with the current state of the target resource.
//	//
//	// In WebDAV remote web authoring, 409 conflict responses are errors sent to the
//	// client so that a user might be able to resolve a conflict and resubmit the
//	// request. [...] Additionally, you may get a 409 response when uploading a file
//	// that is older than the existing one on the server, resulting in a version
//	// control conflict.
//	//
//	// In other systems, 409 responses may be used for implementation-specific purposes,
//	// such as to indicate that the server has received multiple requests to update
//	// the same resource."
//	//
//	// - Quoted from "409 Conflict" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/409
func Conflict(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusConflict),
		WithCode("Conflict"),
		WithMessage("Request conflicts with the current state of the server."),
		WithError(errors.New("user agent sent a request which conflicts with the current state")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// Gone creates a new [Exception] with the "410 Gone" status code, a human
// readable message and error. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 410 Gone client error response status code indicates that the target
//	// resource is no longer available at the origin server and that this condition
//	// is likely to be permanent. A 410 response is cacheable by default.
//	//
//	// Clients should not repeat requests for resources that return a 410 response,
//	// and website owners should remove or replace links that return this code.
//	// If server owners don't know whether this condition is temporary or permanent, Exception
//	// a 404 status code should be used instead."
//	//
//	// - Quoted from "410 Gone" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/410
func Gone(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusGone),
		WithCode("Gone"),
		WithMessage("The requested content has been permanently deleted."),
		WithError(errors.New("user agent has requested content that has been permanently deleted")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// LengthRequired creates a new [Exception] with the "411 Length Required" status
// code, a human readable message and error. The severity of this Exception by
// default is [WARN].
//
//	// "The HTTP 411 Length Required client error response status code indicates that
//	// the server refused to accept the request without a defined Content-Length header."
//	//
//	// - Quoted from "411 Length Required" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/411
func LengthRequired(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusLengthRequired),
		WithCode("Length Required"),
		WithMessage(`The request does not contain a "Content-Length" header, which is required.`),
		WithError(errors.New(`user agent has requested without required "Content-Length" header`)),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// PreconditionFailed creates a new [Exception] with the "412 Precondition Failed"
// status code, a human readable message and error. The severity of this Exception
// by default is [WARN].
//
//	// "The HTTP 412 Precondition Failed client error response status code indicates
//	// that access to the target resource was denied. This happens with conditional
//	// requests on methods other than GET or HEAD when the condition defined by the
//	// If-Unmodified-Since or If-Match headers is not fulfilled. In that case, the
//	// request (usually an upload or a modification of a resource) cannot be made
//	// and this error response is sent back."
//	//
//	// - Quoted from "412 Precondition Failed" by Mozilla Contributors, licensed
//	// under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/412
func PreconditionFailed(err error, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusPreconditionFailed),
		WithCode("Precondition Failed"),
		WithMessage("The request does passes required preconditions."),
		WithError(err),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// ContentTooLarge creates a new [Exception] with the "413 Content Too Large"
// status code, a human readable message and error. The severity of this
// Exception by default is [WARN].
//
//	// "The HTTP 413 Content Too Large client error response status code indicates
//	// that the request entity was larger than limits defined by server. The server
//	// might close the connection or return a Retry-After header field."
//	//
//	// - Quoted from "413 Content Too Large" by Mozilla Contributors, licensed
//	// under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/413
func ContentTooLarge(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusRequestEntityTooLarge),
		WithCode("Content Too Large"),
		WithMessage("The request body is larger than expected."),
		WithError(errors.New("user agent sent request body that is larger than expected")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// URITooLong creates a new [Exception] with the "414 URI Too Long"
// status code, a human readable message and error. The severity of this
// Exception by default is [WARN].
//
//	// "The HTTP 414 URI Too Long client error response status code indicates
//	// that a URI requested by the client was longer than the server is willing
//	// to interpret.
//	//
//	// There are a few rare conditions when this error might occur:
//	//
//	//   - a client has improperly converted a POST request to a GET request with
//	//     long query information,
//	//
//	//   - a client has descended into a loop of redirection (for example, a
//	//     redirected URI prefix that points to a suffix of itself), or
//	//
//	//   - the server is under attack by a client attempting to exploit potential
//	//     security holes."
//	//
//	// - Quoted from "414 URI Too Long" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/414
func URITooLong(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusRequestURITooLong),
		WithCode("URI Too Long"),
		WithMessage("The request has a URI longer than expected."),
		WithError(errors.New("user agent sent request with URI longer than expected")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// UnsupportedMediaType creates a new [Exception] with the "415 Unsupported Media Type"
// status code, a human readable message and error. The severity of this Exception by
// default is [WARN].
//
//	// "The HTTP 415 Unsupported Media Type client error response status code indicates
//	// that the server refused to accept the request because the message content format
//	// is not supported.
//	//
//	// The format problem might be due to the request's indicated Content-Type or
//	// Content-Encoding, or as a result of processing the request message content. Some
//	// servers may be strict about the expected Content-Type of requests. For example,
//	// sending UTF8 instead of UTF-8 to specify the UTF-8 charset may cause the server
//	// to consider the media type invalid."
//	//
//	// - Quoted from "415 Unsupported Media Type" by Mozilla Contributors, licensed
//	// under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/415
func UnsupportedMediaType(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusUnsupportedMediaType),
		WithCode("Unsupported Media Type"),
		WithMessage("The request unsupported media type."),
		WithError(errors.New("user agent sent request with unsupported media type")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// RangeNotSatisfiable creates a new [Exception] with the "416 Range Not Satisfiable"
// status code, a human readable message and error. The severity of this Exception by
// default is [WARN]. A "Content-Range" header is sent with the provided number of
// bytes via the "contentRange" parameter.
//
//	// "The HTTP 416 Range Not Satisfiable client error response status code indicates
//	// that a server could not serve the requested ranges. The most likely reason for
//	// this response is that the document doesn't contain such ranges, or that the
//	// Range header value, though syntactically correct, doesn't make sense.
//	//
//	// The 416 response message should contain a Content-Range indicating an unsatisfied
//	// range (that is a '*') followed by a '/' and the current length of the resource,
//	// e.g., Content-Range: bytes */12777
//	//
//	// When encountering this error, browsers typically either abort the operation
//	// (for example, a download will be considered non-resumable) or request the whole
//	// document again without ranges."
//	//
//	// - Quoted from "416 Range Not Satisfiable" by Mozilla Contributors, licensed
//	// under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/416
func RangeNotSatisfiable(contentRange int, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusUnsupportedMediaType),
		WithCode("Range Not Satisfiable"),
		WithMessage(`Request's "Range" header cannot be satified.`),
		WithError(errors.New(`user agent sent request with unsitisfiable "Range" header`)),
		WithSeverity(WARN),

		WithHeader("Content-Range", fmt.Sprintf("bytes */%d", contentRange)),
	}
	o = append(o, opts...)

	return newException(o...)
}

// ExpectationFailed creates a new [Exception] with the "417 Expectation Failed"
// status code, a human readable message and error. The severity of this Exception
// by default is [WARN].
//
//	// "The HTTP 417 Expectation Failed client error response status code indicates
//	// that the expectation given in the request's Expect header could not be met.
//	// After receiving a 417 response, a client should repeat the request without an
//	// Expect request header, including the file in the request body without waiting
//	// for a 100 response. See the Expect header documentation for more details."
//	//
//	// - Quoted from "417 Expectation Failed" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/417
func ExpectationFailed(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusExpectationFailed),
		WithCode("Exception Failed"),
		WithMessage(`Request's "Expect" header cannot be met.`),
		WithError(errors.New(`user agent sent request with unmetable "Expect" header`)),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// ImATeapot creates a new [Exception] with the "418 I'm a teapot" status code,
// a human readable message and error. The severity of this Exception by default
// is [WARN].
//
//	// "The HTTP 418 I'm a teapot status response code indicates that the server
//	// refuses to brew coffee because it is, permanently, a teapot. A combined
//	// coffee/tea pot that is temporarily out of coffee should instead return 503.
//	// This error is a reference to Hyper Text Coffee Pot Control Protocol defined
//	// in April Fools' jokes in 1998 and 2014.
//	//
//	// Some websites use this response for requests they do not wish to handle,
//	// such as automated queries."
//	//
//	// - Quoted from "418 I'm a teapot" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/418
func ImATeapot(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusTeapot),
		WithCode("I'm a teapot"),
		WithMessage("The request to brew coffee with a teapot is impossible."),
		WithError(errors.New("user agent tried to brew coffee with a teapot, an impossible request")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// MisdirectedRequest creates a new [Exception] with the "421 Misdirected Request"
// status code, a human readable message and error. The severity of this Exception
// by default is [WARN].
//
//	// "The HTTP 421 Misdirected Request client error response status code indicates
//	// that the request was directed to a server that is not able to produce a response.
//	// This can be sent by a server that is not configured to produce responses for the
//	// combination of scheme and authority that are included in the request URI.
//	//
//	// Clients may retry the request over a different connection."
//	//
//	// - Quoted from "421 Misdirected Request" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/421
func MisdirectedRequest(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusMisdirectedRequest),
		WithCode("Misdirected Request"),
		WithMessage("The request was directed to a location unable to produce a response."),
		WithError(errors.New("user agent requested on a location that is unable to produce a response.")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// UnprocessableContent creates a new [Exception] with the "422 Unprocessable Content"
// status code, a human readable message and error. The severity of this Exception
// by default is [WARN].
//
//	// "The HTTP 422 Unprocessable Content client error response status code indicates
//	// that the server understood the content type of the request content, and the syntax
//	// of the request content was correct, but it was unable to process the contained
//	// instructions.
//	//
//	// Clients that receive a 422 response should expect that repeating the request
//	// without modification will fail with the same error."
//	//
//	// - Quoted from "422 Unprocessable Content" by Mozilla Contributors, licensed
//	// under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/422
func UnprocessableContent(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusUnprocessableEntity),
		WithCode("Unprocessable Content"),
		WithMessage("Unable to follow request due to semantic errors."),
		WithError(errors.New("user agent sent request containing semantic errors.")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// Locked creates a new [Exception] with the "423 Locked" status code, a human
// readable message and error. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 423 Locked client error response status code indicates that a
//	// resource is locked, meaning it can't be accessed. Its response body should
//	// contain information in WebDAV's XML format.
//	//
//	// NOTE: The ability to lock a resource to prevent conflicts is specific to
//	// some WebDAV servers. Browsers accessing web pages will never encounter
//	// this status code; in the erroneous cases it happens, they will handle
//	// it as a generic 400 status code."
//	//
//	// - Quoted from "423 Locked" by Mozilla Contributors, licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/423
func Locked(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusUnprocessableEntity),
		WithCode("Locked"),
		WithMessage("This resource is locked."),
		WithError(errors.New("user agent requested a locked resource")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// FailedDependency creates a new [Exception] with the "424 Failed Dependency"
// status code, a human readable message and error. The severity of this
// Exception by default is [WARN].
//
//	// "The HTTP 424 Failed Dependency client error response status code
//	// indicates that the method could not be performed on the resource
//	// because the requested action depended on another action, and that
//	// action failed.
//	//
//	// Regular web servers typically do not return this status code, but
//	// some protocols like WebDAV can return it. For example, in WebDAV,
//	// if a PROPPATCH request was issued, and one command fails then
//	// automatically every other command will also fail with
//	// 424 Failed Dependency."
//	//
//	// - Quoted from "424 Failed Dependency" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/424
func FailedDependency(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusFailedDependency),
		WithCode("Failed Dependency"),
		WithMessage("Cannot respond due to failure of a previous request."),
		WithError(errors.New("request cannot be due to failure of a previous request")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// TooEarly creates a new [Exception] with the "425 Too Early" status code,
// a human readable message and error. The severity of this Exception by
// default is [WARN].
//
//	// "LIMITED AVAILABILITY
//	//
//	// The HTTP 425 Too Early client error response status code indicates
//	// that the server was unwilling to risk processing a request that might
//	// be replayed to avoid potential replay attacks.
//	//
//	// If a client has interacted with a server recently, early data
//	// (also known as zero round-trip time (0-RTT) data) allows the client
//	// to send data to a server in the first round trip of a connection,
//	// without waiting for the TLS handshake to complete. A client that
//	// sends a request in early data does not need to include the
//	// Early-Data header."
//	//
//	// - Quoted from "425 Too Early" by Mozilla Contributors, licensed under
//	// CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/425
func TooEarly(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusTooEarly),
		WithCode("Too Early"),
		WithMessage("Unwilling to process the request to avoid replay attacks."),
		WithError(errors.New("request was rejected to avoid replay attacks")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// UpgradeRequired creates a new [Exception] with the "426 Upgrade Required"
// status code, a human readable message and error. The severity of this
// Exception by default is [WARN]. A "Upgrade" header is sent with the value
// provided by the "upgrade" parameter.
//
//	// "The HTTP 426 Upgrade Required client error response status code
//	// indicates that the server refused to perform the request using the
//	// current protocol but might be willing to do so after the client
//	// upgrades to a different protocol.
//	//
//	// The server sends an Upgrade header with this response to indicate
//	// the required protocol(s)."
//	//
//	// - Quoted from "426 Upgrade Required" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/426
func UpgradeRequired(upgrade string, opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusUpgradeRequired),
		WithCode("Upgrade Required"),
		WithMessage("An upgrade to the protocol is required to do this request."),
		WithError(fmt.Errorf("user agent needs to upgrade to protocol %q", upgrade)),
		WithHeader("Upgrade", upgrade),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// PreconditionRequired creates a new [Exception] with the "428 Precondition Required"
// status code, a human readable message and error. The severity of this Exception
// by default is [WARN].
//
//	// "The HTTP 428 Precondition Required client error response status code indicates
//	// that the server requires the request to be conditional.
//	//
//	// Typically, a 428 response means that a required precondition header such as
//	// If-Match is missing. When a precondition header does not match the server-side
//	// state, the response should be 412 Precondition Failed."
//	//
//	// - Quoted from "428 Precondition Required" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/426
func PreconditionRequired(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusPreconditionRequired),
		WithCode("Precondition Required"),
		WithMessage("The request needs to be conditional."),
		WithError(errors.New("user agent sent request that is not conditional")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// TooManyRequests creates a new [Exception] with the "429 Too Many Requests"
// status code, a human readable message and error. The severity of this
// Exception by default is [WARN].
//
// To provide a "Retry-After" header, use the [WithHeader] option function.
//
//	// "The HTTP 429 Too Many Requests client error response status code
//	// indicates the client has sent too many requests in a given amount
//	// of time. This mechanism of asking the client to slow down the rate
//	// of requests is commonly called "rate limiting".
//	//
//	// A Retry-After header may be included to this response to indicate
//	// how long a client should wait before making the request again."
//	//
//	// - Quoted from "429 Too Many Requests" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/429
func TooManyRequests(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusTooManyRequests),
		WithCode("Too Many Requests"),
		WithMessage("Too many requests were sent in the span of a short time."),
		WithError(errors.New("user agent sent too many requests")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// RequestHeaderFieldsTooLarge creates a new [Exception] with the
// "431 Request Header Fields Too Large" status code, a human readable
// message and error. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 431 Request Header Fields Too Large client error response
//	// status code indicates that the server refuses to process the request
//	// because the request's HTTP headers are too long. The request may
//	// be resubmitted after reducing the size of the request headers.
//	//
//	// 431 can be used when the total size of request headers is too large
//	// or when a single header field is too large. To help clients running
//	// into this error, indicate which of the two is the problem in the
//	// response body and, ideally, say which headers are too large. This
//	// lets people attempt to fix the problem, such as by clearing cookies.
//	//
//	// Servers will often produce this status if:
//	// - The Referer URL is too long
//	// - There are too many Cookies in the request"
//	//
//	// - Quoted from "431 Request Header Fields Too Large" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/431
func RequestHeaderFieldsTooLarge(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusRequestHeaderFieldsTooLarge),
		WithCode("Request Header Fields Too Large"),
		WithMessage("Headers fields are larger than expected."),
		WithError(errors.New("user agent sent header fields larger than expected")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}

// UnavailableForLegalReasons creates a new [Exception] with the
// "451 Unavailable For Legal Reasons" status code, a human readable
// message and error. The severity of this Exception by default is [WARN].
//
//	// "The HTTP 451 Unavailable For Legal Reasons client error response
//	// status code indicates that the user requested a resource that is
//	// not available due to legal reasons, such as a web page for which
//	// a legal action has been issued."
//	//
//	// - Quoted from "451 Unavailable For Legal Reasons" by Mozilla Contributors,
//	// licensed under CC-BY-SA 2.5.
//	//
//	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/451
func UnavailableForLegalReasons(opts ...Option) Exception {
	o := []Option{
		WithStatus(http.StatusUnavailableForLegalReasons),
		WithCode("Unavailable For Legal Reasons"),
		WithMessage("Content cannot be legally be provided."),
		WithError(errors.New("user agent requested content that cannot be legally provided")),
		WithSeverity(WARN),
	}
	o = append(o, opts...)

	return newException(o...)
}
