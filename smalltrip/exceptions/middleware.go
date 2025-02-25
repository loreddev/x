package exceptions

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func HandlerJSON(fallback HandlerFunc) HandlerFunc {
	return func(e Exception, w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(e)
		if err != nil {
			e.Err = errors.Join(fmt.Errorf("marshalling Exception struct: %s", e.Error()), e.Err)

			fallback(e, w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(e.Status)

		_, err = w.Write(j)
		if err != nil {
			e.Err = errors.Join(fmt.Errorf("writing JSON response to body: %s", e.Error()), e.Err)

			HandlerText(e, w, r)
			return
		}
	}
}

var _ HandlerFunc = HandlerJSON(HandlerText)

func HandlerText(e Exception, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(e.Status)

	_, err := w.Write([]byte(fmt.Sprintf(
		"Status: %3d\n"+
			"Code: %s"+
			"Message: %s\n"+
			"Err: %s\n"+
			"Severity: %s\n\n"+
			"%+v\n\n"+
			"%#v",
		e.Status,
		e.Code,
		e.Message,
		e.Err,
		e.Severity,

		e, e,
	)))
	if err != nil {
		_, _ = w.Write([]byte(fmt.Sprintf(
			"Ok, what should we do at this point? You fucked up so bad that this message " +
				"shouldn't even be able to be sent in the first place. If you are a normal user I'm " +
				"so sorry for you to be reading this. If you're a developer, go fix your ResponseWriter " +
				"implementation, because this should never happen in any normal codebase. " +
				"I hope for the life of anyone you love you don't use this message in some " +
				"error checking or any sort of API-contract, because there will be no more hope " +
				"for you or your project. May God or any other or any other divinity that you may " +
				"or may not believe be with you when trying to fix this mistake, you will need it.",
			// If someone use this as part of the API-contract I'll not even be surprised.
			// So any change to this message is still considered a breaking change.
		)))
	}
}

var _ HandlerFunc = HandlerText
