package exceptions

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

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
