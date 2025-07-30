package problem

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"slices"
	"text/template"
)

type Problem interface {
	Type() string
	Title() string
	Status() int
	Detail() string
	Instance() string

	Handler(self Problem) http.Handler

	http.Handler
}

type RegisteredProblem struct {
	TypeURI       string `json:"type,omitempty"     xml:"type,omitempty"`
	TypeTitle     string `json:"title,omitempty"    xml:"title,omitempty"`
	StatusCode    int    `json:"status,omitempty"   xml:"status,omitempty"`
	DetailMessage string `json:"detail,omitempty"   xml:"detail,omitempty"`
	InstanceURI   string `json:"instance,omitempty" xml:"instance,omitempty"`

	XMLName xml.Name `json:"-" xml:"problem"`

	handler Handler `json:"-" xml:"-"`
}

func NewStatus(s int, opts ...Option) RegisteredProblem {
	return New(slices.Concat([]Option{WithStatus(s)}, opts)...)
}

func NewDetailed(s int, detail string, opts ...Option) RegisteredProblem {
	return New(slices.Concat([]Option{WithStatus(s), WithDetail(detail)}, opts)...)
}

func New(opts ...Option) RegisteredProblem {
	p := RegisteredProblem{
		TypeURI: DefaultTypeURI,
		handler: DefaultHandler,
	}

	for _, opt := range opts {
		opt(&p)
	}
	return p
}

var (
	DefaultTypeURI  = "about:blank"
	DefaultTemplate = template.Must(template.New("x-smalltrip-problem-template").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<title>{{ .Status }} - {{ .Title }}</title>
	</head>
	<body>
		<h1>{{.Status}} - {{ .Title }}</h1>
		<p><code>{{ .Type }}</code></p>
		<p>{{ .Detail }}</p>
		{{if .Instance}}
			<p>Instance: {{ .Instance }}</p>
		{{end}}
		<code>{{printf "%#v" .}}<code>
	</body>
<html>
`))
	DefaultHandler = HandlerMiddleware(HandlerBrowser(DefaultTemplate))
)

func (p RegisteredProblem) Type() string {
	return p.TypeURI
}

func (p RegisteredProblem) Title() string {
	return p.TypeTitle
}

func (p RegisteredProblem) Status() int {
	return p.StatusCode
}

func (p RegisteredProblem) Detail() string {
	return p.DetailMessage
}

func (p RegisteredProblem) Instance() string {
	return p.InstanceURI
}

func (p RegisteredProblem) Handler(self Problem) http.Handler {
	return p.handler(self)
}

func (p RegisteredProblem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.Handler(p).ServeHTTP(w, r)
}

func WithType(t string) Option {
	return func(p *RegisteredProblem) {
		p.TypeURI = t
	}
}

func WithTitle(t string) Option {
	return func(p *RegisteredProblem) {
		p.TypeTitle = t
	}
}

func WithStatus(s int) Option {
	return func(p *RegisteredProblem) {
		if p.TypeTitle == "" {
			p.TypeTitle = http.StatusText(s)
		}
		p.StatusCode = s
	}
}

func WithDetail(d string) Option {
	return func(p *RegisteredProblem) {
		p.DetailMessage = d
	}
}

func WithDetailf(f string, args ...any) Option {
	return func(p *RegisteredProblem) {
		p.DetailMessage = fmt.Sprintf(f, args...)
	}
}

func WithError(err error) Option {
	return func(p *RegisteredProblem) {
		p.DetailMessage = err.Error()
	}
}

func WithInstance(i string) Option {
	return func(p *RegisteredProblem) {
		p.InstanceURI = i
	}
}

type Option func(*RegisteredProblem)
