package rerrors

import (
	"net/http"
	"strconv"
)

func BadRequest(reason ...string) RouteError {
	info := map[string]any{}

	if len(reason) == 1 {
		info["reason"] = reason[0]
	} else if len(reason) > 1 {
		for i, r := range reason {
			info["reason_"+strconv.Itoa(i)] = r
		}
	}

	return NewRouteError(http.StatusBadRequest, "Bad Request", info)
}

func NotFound() RouteError {
	return NewRouteError(http.StatusNotFound, "Not Found", map[string]any{})
}

func MissingCookies(cookies []string) RouteError {
	return NewRouteError(http.StatusBadRequest, "Missing cookies", map[string]any{
		"missing_cookies": cookies,
	})
}

func MethodNowAllowed(method string, allowedMethods []string) RouteError {
	return NewRouteError(http.StatusMethodNotAllowed, "Method not allowed", map[string]any{
		"method":          method,
		"allowed_methods": allowedMethods,
	})
}

func MissingParameters(params []string) RouteError {
	return NewRouteError(http.StatusBadRequest, "Missing parameters", map[string]any{
		"missing_parameters": params,
	})
}
