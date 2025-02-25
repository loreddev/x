package exceptions

import (
	"errors"
	"fmt"
	"net/http"
)

type Exception struct {
	Status   int      `json:"status"`          // HTTP Status Code
	Code     string   `json:"code"`            // Application error code
	Message  string   `json:"message"`         // User friendly message
	Err      error    `json:"error,omitempty"` // Go error
	Severity Severity `json:"severity"`        // Exception level
}

var (
	_ fmt.Stringer = Exception{}
	_ error        = Exception{}
	_ http.Handler = Exception{}
)

func (e Exception) String() string {
	return fmt.Sprintf("%s %3d %s Exception %q", e.Severity, e.Status, e.Code, e.Message)
}

func (e Exception) Error() string {
	return e.String()
}

func (e Exception) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e.handler != nil {
		e.handler(e, w, r)
	}

	handler, ok := r.Context().Value(handlerFuncCtxKey).(HandlerFunc)
	if !ok {
		e.handler(e, w, r)
	}

	handler(e, w, r)
}
