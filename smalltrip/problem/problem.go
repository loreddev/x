package problem

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"slices"
)

type Problem struct {
	Type     string `json:"type,omitempty"     xml:"type,omitempty"`
	Title    string `json:"title,omitempty"    xml:"title,omitempty"`
	Status   int    `json:"status,omitempty"   xml:"status,omitempty"`
	Detail   string `json:"detail,omitempty"   xml:"detail,omitempty"`
	Instance string `json:"instance,omitempty" xml:"instance,omitempty"`

	XMLName xml.Name `json:"-" xml:"problem"`

	handler func(any) http.Handler `json:"-" xml:"-"`
}

func New(opts ...Option) Problem {
	p := Problem{
		Type:    DefaultType,
		handler: ProblemHandler,
	}
	for _, opt := range opts {
		opt(&p)
	}
	return p
}

const DefaultType = "about:blank"

func NewStatus(s int, opts ...Option) Problem {
	return New(slices.Concat([]Option{WithStatus(s)}, opts)...)
}

func NewDetailed(s int, detail string, opts ...Option) Problem {
	return New(slices.Concat([]Option{WithStatus(s), WithDetail(detail)}, opts)...)
}

func (p Problem) StatusCode() int {
	return p.Status
}

func (p Problem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.handler(p).ServeHTTP(w, r)
}

func WithType(t string) Option {
	return func(p *Problem) {
		p.Type = t
	}
}

func WithTitle(t string) Option {
	return func(p *Problem) {
		p.Title = t
	}
}

func WithStatus(s int) Option {
	return func(p *Problem) {
		if p.Title == "" {
			p.Title = http.StatusText(s)
		}

		p.Status = s
	}
}

func WithDetail(d string) Option {
	return func(p *Problem) {
		p.Detail = d
	}
}

func WithDetailf(f string, args ...any) Option {
	return func(p *Problem) {
		p.Detail = fmt.Sprintf(f, args...)
	}
}

func WithError(err error) Option {
	return func(p *Problem) {
		p.Detail = err.Error()
	}
}

func WithInstance(i string) Option {
	return func(p *Problem) {
		p.Instance = i
	}
}

type Option func(*Problem)
