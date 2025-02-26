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
