package problem

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

type Handler func(p Problem) http.Handler

func HandlerAll(p Problem) http.Handler {
	h := HandlerContentType(map[string]Handler{
		"application/xml":    HandlerXML,
		ProblemMediaTypeXML:  HandlerXML,
		"application/json":   HandlerJSON,
		ProblemMediaTypeJSON: HandlerJSON,
	}, HandlerJSON)
	return h(p)
}

func HandlerContentType(handlers map[string]Handler, fallback ...Handler) Handler {
	return func(p Problem) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for t, h := range handlers {
				if strings.Contains(r.Header.Get("Accept"), t) {
					h(p).ServeHTTP(w, r)
					return
				}
			}
			if len(fallback) > 0 {
				fallback[0](p).ServeHTTP(w, r)
			}
		})
	}
}

func HandlerXML(p Problem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ProblemMediaTypeXML)

		b, err := xml.Marshal(p)
		if err != nil {
			HandlerJSON(p).ServeHTTP(w, r)
			return
		}

		w.WriteHeader(p.Status())

		_, err = w.Write(b)
		if err != nil {
			HandlerJSON(p).ServeHTTP(w, r)
		}
	})
}

func HandlerJSON(p Problem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ProblemMediaTypeJSON)

		b, err := json.Marshal(p)
		if err != nil {
			HandlerText(p).ServeHTTP(w, r)
			return
		}

		w.WriteHeader(p.Status())

		_, err = w.Write(b)
		if err != nil {
			HandlerText(p).ServeHTTP(w, r)
		}
	})
}

func HandlerText(p Problem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		w.WriteHeader(p.Status())

		s := fmt.Sprintf(
			"Type: %s\n"+
				"Status: %3d\n"+
				"Title: %s\n"+
				"Detail: %s\n"+
				"Instance: %s\n\n"+
				p.Type(),
			p.Status(),
			p.Title(),
			p.Detail(),
			p.Instance(),
		)

		_, err := w.Write(fmt.Appendf([]byte{}, "%s%+v\n\n%#v", s, p, p))
		if err != nil {
			_, _ = w.Write(fmt.Append([]byte{},
				"Ok, what should we do at this point? You fucked up so bad that this message "+
					"shouldn't even be able to be sent in the first place. If you are a normal user I'm "+
					"so sorry for you to be reading this. If you're a developer, go fix your ResponseWriter "+
					"implementation, because this should never happen in any normal codebase. "+
					"I hope for the life of anyone you love you don't use this message in some "+
					"error checking or any sort of API-contract, because there will be no more hope "+
					"for you or your project. May God or any other divinity that you may "+
					"or may not believe be with you when trying to fix this mistake, you will need it.",
				// If someone use this as part of the API-contract I'll not even be surprised.
				// So any change to this message is still considered a breaking change.
			))
		}
	})
}

const (
	ProblemMediaTypeJSON = "application/problem+json"
	ProblemMediaTypeXML  = "application/problem+xml"
)
