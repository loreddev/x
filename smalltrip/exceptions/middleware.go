package exceptions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"strings"

	"forge.capytal.company/loreddev/x/groute/middleware"
)

func Middleware(options ...MiddlewareOption) middleware.Middleware {
	opts := middlewareOpts{
		templates:      make(map[int]*template.Template),
		handlers:       make(map[string]HandlerFunc),
		defaultHandler: HandlerJSON(HandlerText),
	}

	for _, option := range options {
		option(&opts)
	}

	if _, ok := opts.templates[0]; !ok {
		opts.templates[0] = defaultTemplate
	}

	if _, ok := opts.handlers["application/json"]; !ok {
		opts.handlers["application/json"] = HandlerJSON(HandlerText)
	}
	if _, ok := opts.handlers["text/html"]; !ok {
		opts.handlers["text/html"] = HandlerTemplates(opts.templates, opts.defaultHandler)
	}
	if _, ok := opts.handlers["application/xhtml+xml"]; !ok {
		opts.handlers["application/xhtml+xml"] = HandlerTemplates(opts.templates, opts.defaultHandler)
	}
	if _, ok := opts.handlers["application/xml"]; !ok {
		opts.handlers["application/xml"] = HandlerTemplates(opts.templates, opts.defaultHandler)
	}

	return NewMiddleware(func(e Exception, w http.ResponseWriter, r *http.Request) {
		for k, v := range opts.handlers {
			if strings.Contains(r.Header.Get("Accept"), k) {
				v(e, w, r)
				return
			}
		}
		opts.defaultHandler(e, w, r)
	})
}

var defaultTemplate = template.Must(template.New("xx-small-trip-default-Exception-template").Parse(`
Status: {{ .Status }}
Code: {{ .Code }}
Message: {{ .Message }}
Err: {{ .Err }}
Severity: {{ .Severity }}
`))

type MiddlewareOption = func(*middlewareOpts)

func MiddlewareTemplate(t *template.Template, statusCode ...int) MiddlewareOption {
	return func(mo *middlewareOpts) {
		if len(statusCode) > 0 {
			mo.templates[statusCode[0]] = t
		} else {
			mo.templates[0] = t
		}
	}
}

func MiddlewareHandler(h HandlerFunc, mimeType ...string) MiddlewareOption {
	return func(mo *middlewareOpts) {
		if len(mimeType) > 0 {
			mo.handlers[mimeType[0]] = h
		} else {
			mo.defaultHandler = h
		}
	}
}

type middlewareOpts struct {
	templates      map[int]*template.Template
	handlers       map[string]HandlerFunc
	defaultHandler HandlerFunc
}

func NewMiddleware(handler HandlerFunc) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), handlerFuncCtxKey, handler))
			next.ServeHTTP(w, r)
		})
	}
}

const handlerFuncCtxKey = "xx-smalltrip-Exception-handler-func"

type HandlerFunc = func(e Exception, w http.ResponseWriter, r *http.Request)

func HandlerTemplates(ts map[int]*template.Template, fallback HandlerFunc) HandlerFunc {
	return func(e Exception, w http.ResponseWriter, r *http.Request) {
		if len(ts) == 0 {
			fallback(e, w, r)
			return
		}

		t, ok := ts[e.Status]
		if ok {
			HandlerTemplate(t, fallback)(e, w, r)
			return
		}

		// Loops over ordered list and gets the last one that is small or equal
		// the current Status. For example, if the current Exception has Status 404,
		// and we provide a map like:
		//
		// map[int]*template.Template{
		//   100: Template100,
		//   200: Template200,
		//   300: Template300,
		//   400: Template400,
		//   500: Template500,
		// }
		//
		// It will be converted to a ordered list of keys: 100, 200, 300, 400, 500.
		// This loops iterates on all keys until a value bigger than the current Status
		// is found (in this example, 500), then it uses the previous (in this example 400).
		//
		// So the 404 Exception will be rendered using the Template400.

		keys := make([]int, len(ts), len(ts))

		var i int
		for k := range ts {
			keys[i] = k
			i++
		}

		slices.Sort(keys)

		key := keys[0]
		for _, k := range keys {
			if k > e.Status {
				break
			}
			key = k
		}

		t, ok = ts[key]
		if ok {
			HandlerTemplate(t, fallback)(e, w, r)
			return
		}

		fallback(e, w, r)
	}
}

func HandlerTemplate(t *template.Template, fallback HandlerFunc) HandlerFunc {
	return func(e Exception, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		for k := range e.headers {
			w.Header().Set(k, e.headers.Get(k))
		}

		w.WriteHeader(e.Status)

		err := t.Execute(w, e)
		if err != nil {
			e.Err = errors.Join(fmt.Errorf("executing Exception template: %s", e.Error()), e.Err)

			fallback(e, w, r)
			return
		}
	}
}

func HandlerJSON(fallback HandlerFunc) HandlerFunc {
	return func(e Exception, w http.ResponseWriter, r *http.Request) {
		j, err := json.Marshal(e)
		if err != nil {
			e.Err = errors.Join(fmt.Errorf("marshalling Exception struct: %s", e.Error()), e.Err)

			fallback(e, w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		for k := range e.headers {
			w.Header().Set(k, e.headers.Get(k))
		}

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
	for k := range e.headers {
		w.Header().Set(k, e.headers.Get(k))
	}

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
