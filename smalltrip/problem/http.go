package problem

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

func ProblemHandler(p any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Accept")
		if strings.Contains(h, "application/xml") || strings.Contains(h, ProblemMediaTypeXML) {
			ProblemHandlerXML(p).ServeHTTP(w, r)
			return
		}
		ProblemHandlerJSON(p).ServeHTTP(w, r)
	})
}

func ProblemHandlerXML(p any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// w.Header().Set("Content-Type", ProblemMediaTypeXML)
		w.Header().Set("Content-Type", "application/xml")

		b, err := xml.Marshal(p)
		if err != nil {
			ProblemHandlerJSON(p).ServeHTTP(w, r)
			return
		}

		w.WriteHeader(GetStatus(p))

		_, err = w.Write(b)
		if err != nil {
			ProblemHandlerJSON(p).ServeHTTP(w, r)
		}
	})
}

func ProblemHandlerJSON(p any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", ProblemMediaTypeJSON)

		b, err := json.Marshal(p)
		if err != nil {
			ProblemHandlerText(p).ServeHTTP(w, r)
			return
		}

		w.WriteHeader(GetStatus(p))

		_, err = w.Write(b)
		if err != nil {
			ProblemHandlerText(p).ServeHTTP(w, r)
		}
	})
}

func ProblemHandlerText(p any) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		w.WriteHeader(GetStatus(p))

		var s string
		if p, ok := p.(Problem); ok {
			s = fmt.Sprintf(
				"Type: %s\n"+
					"Status: %3d\n"+
					"Title: %s\n"+
					"Detail: %s\n"+
					"Instance: %s\n\n"+
					p.Type,
				p.Status,
				p.Title,
				p.Detail,
				p.Instance,
			)
		}
		if p, ok := p.(*Problem); ok {
			s = fmt.Sprintf(
				"Type: %s\n"+
					"Status: %3d\n"+
					"Title: %s\n"+
					"Detail: %s\n"+
					"Instance: %s\n\n"+
					p.Type,
				p.Status,
				p.Title,
				p.Detail,
				p.Instance,
			)
		}

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

func GetStatus(p any) int {
	if p, ok := p.(interface{ StatusCode() int }); ok {
		return p.StatusCode()
	}
	return http.StatusInternalServerError
}

const (
	ProblemMediaTypeJSON = "application/problem+json"
	ProblemMediaTypeXML  = "application/problem+xml"
)
