package problem

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"
)

func ServeProblem(p Problem, w http.ResponseWriter, r *http.Request) {
	acc := r.Header.Get("Accept")
	if strings.Contains(acc, "application/problem+xml") || strings.Contains(acc, "application/xml") {
		ServeProblemXML(p, w, r)
		return
	}
	ServeProblemJSON(p, w, r)
}

func ServeProblemXML(p Problem, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml") // TODO: Change to problem+xml

	w.WriteHeader(p.Status())

	var b []byte
	var err error

	if h := r.Header.Get("User-Agent"); strings.Contains(h, "Mozilla") ||
		strings.Contains(h, "WebKit") ||
		strings.Contains(h, "Chrome") {
		b, err = xml.MarshalIndent(p, "", "  ")
	} else {
		b, err = xml.Marshal(p)
	}

	if err != nil {
		ServeProblemJSON(p, w, r)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		ServeProblemJSON(p, w, r)
	}
}

func ServeProblemJSON(p Problem, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/problem+json")

	w.WriteHeader(p.Status())

	var b []byte
	var err error

	if h := r.Header.Get("User-Agent"); strings.Contains(h, "Mozilla") ||
		strings.Contains(h, "WebKit") ||
		strings.Contains(h, "Chrome") {
		b, err = json.MarshalIndent(p, "", "  ")
	} else {
		b, err = json.Marshal(p)
	}

	if err != nil {
		ServeProblemText(p, w, r)
		return
	}

	_, err = w.Write(b)
	if err != nil {
		ServeProblemText(p, w, r)
	}
}

func ServeProblemText(p Problem, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/problem+text")

	s := ""

	w.WriteHeader(p.Status())

	if p, ok := p.(interface{ Type() string }); ok {
		s += fmt.Sprintf("Type: %s\n", p.Type())
	}

	if p, ok := p.(interface{ Title() string }); ok {
		s += fmt.Sprintf("Title: %s\n", p.Title())
	}

	s += fmt.Sprintf("Status: %3d\n", p.Status())

	if p, ok := p.(interface{ Detail() string }); ok {
		s += fmt.Sprintf("Detail: %s\n", p.Detail())
	}

	if p, ok := p.(interface{ Instance() string }); ok {
		s += fmt.Sprintf("Instance: %s\n", p.Instance())
	}

	_, err := w.Write(fmt.Appendf([]byte{}, "%s\n\n%+v\n\n%#v", s, p, p))
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
}
