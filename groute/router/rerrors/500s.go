package rerrors

import (
	"errors"
	"net/http"
)

func InternalError(errs ...error) RouteError {
	err := errors.Join(errs...)
	return NewRouteError(http.StatusInternalServerError, "Internal server error", map[string]any{
		"errors":      err,
		"errors_desc": err.Error(),
	})
}
