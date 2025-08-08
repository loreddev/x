package multiplexer

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
)

func WithPatternsOptions(mux Multiplexer, opts ...PatternOption) Multiplexer {
	return &patternsMux{opts: opts, Multiplexer: mux}
}

func AddTrailingSlash() PatternOption {
	return func(s string) string {
		if strings.HasSuffix(s, "...}") {
			return s
		}
		// If the pattern is /image.html{$}, we modify it to be /image.html/{$}
		if strings.HasSuffix(s, "{$}") && !strings.HasSuffix(s, "/{$}") {
			return strings.TrimSuffix(s, "{$}") + "/{$}"
		}
		if !strings.HasSuffix(s, "/") {
			return s + "/"
		}
		return s
	}
}

func RemoveTrailingSlash() PatternOption {
	return func(s string) string {
		s = strings.TrimSuffix(s, "/")
		if strings.HasSuffix(s, "/{$}") {
			return strings.TrimSuffix(s, "/{$}") + "{$}"
		}
		return s
	}
}

func AddStrictEnd() PatternOption {
	return func(s string) string {
		if strings.HasSuffix(s, "{$}") {
			return s
		}
		return s + "{$}"
	}
}

func RemoteStrictEnd() PatternOption {
	return func(s string) string {
		return strings.TrimSuffix(s, "{$}")
	}
}

type PatternOption func(string) string

func WithPatternRules(mux Multiplexer, rules ...PatternRule) Multiplexer {
	opts := make([]PatternOption, len(rules))
	for i, r := range rules {
		opts[i] = func(s string) string { r(s); return s }
	}
	return &patternsMux{Multiplexer: mux, opts: opts}
}

func NoTrailingSlash() PatternRule {
	return func(s string) {
		if strings.HasSuffix(s, "/") || strings.HasSuffix(s, "/{$}") {
			panic(fmt.Sprintf("no-trailing-slash: pattern %q has trailing slash", s))
		}
	}
}

func EnsureTrailingSlash() PatternRule {
	return func(s string) {
		if !strings.HasSuffix(s, "/{$}") && !strings.HasSuffix(s, "/") && !strings.HasSuffix(s, "...}") {
			panic(fmt.Sprintf("trailing-slash: pattern %q doesn't has a trailing slash", s))
		}
	}
}

func NoMethod() PatternRule {
	return func(s string) {
		if len(strings.Split(s, " ")) > 1 {
			panic(fmt.Sprintf("no-method: pattern %q has a method", s))
		}
	}
}

func EnsureMethod(methods ...string) PatternRule {
	if len(methods) == 0 && methods != nil {
		methods = DefaultMethods
	}
	return func(s string) {
		sp := strings.Split(s, " ")
		if len(sp) <= 0 {
			panic(fmt.Sprintf("method: pattern %q doesn't has a method", s))
		}
		if methods != nil {
			if slices.Contains(methods, sp[0]) {
				panic(fmt.Sprintf("method: pattern %q doesn't has a valid method, valid methods are: %s", s, strings.Join(methods, ", ")))
			}
		}
	}
}

var DefaultMethods = []string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

func EnsureStrictEnd() PatternRule {
	return func(s string) {
		if !strings.HasSuffix(s, "{$}") && !strings.HasSuffix(s, "...}") {
			panic(fmt.Sprintf(`strict-end: pattern %q doesn't end with "{$}"`, s))
		}
	}
}

func NoStrictEnd() PatternRule {
	return func(s string) {
		if strings.HasSuffix(s, "{$}") {
			panic(fmt.Sprintf(`no-strict-end: pattern %q ends with "{$}"`, s))
		}
	}
}

type PatternRule func(string)

type patternsMux struct {
	opts []PatternOption
	Multiplexer
}

func (pm *patternsMux) HandleFunc(p string, h func(http.ResponseWriter, *http.Request)) {
	for _, po := range pm.opts {
		p = po(p)
	}
	pm.Multiplexer.HandleFunc(p, h)
}

func (pm *patternsMux) Handle(p string, h http.Handler) {
	for _, po := range pm.opts {
		p = po(p)
	}
	pm.Multiplexer.Handle(p, h)
}
