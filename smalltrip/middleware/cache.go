// Copyright 2025-present Gustavo "Guz" L. de Mello
// Copyright 2025-present The Lored.dev Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func Cache(options ...CacheOption) Middleware {
	d := defaultCacheDirectives

	for _, option := range options {
		option(&d)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", d.String())
			next.ServeHTTP(w, r)
		})
	}
}

// TODO: SmartCache is a smarter implementation of Cache that handles requests
// with authorization, Cache-Control from the client, and others.
func SmartCache(options ...CacheOption) Middleware {
	return Cache(options...)
}

// TODO: PersistentCache is a smarter implementation of SmartCache that handles requests
// with authorization, Cache-Control from the client, and stores responses into
// a persistent storage solution like Redis.
func PersistentCache(options ...CacheOption) Middleware {
	return SmartCache(options...)
}

type CacheOption func(*directives)

func CacheMaxAge(t time.Duration) CacheOption {
	return func(d *directives) { d.maxAge = &t }
}

func CacheSMaxAge(t time.Duration) CacheOption {
	return func(d *directives) { d.sMaxage = &t }
}

func CacheNoCache(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.noCache = &bool }
}

func CacheNoStore(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.noStore = &bool }
}

func CacheNoTransform(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.noTransform = &bool }
}

func CacheMustRevalidate(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.mustRevalidate = &bool }
}

func CacheProxyRevalidate(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.proxyRevalidate = &bool }
}

func CacheMustUnderstand(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.mustUnderstand = &bool }
}

func CachePrivate(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.private = &bool }
}

func CachePublic(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.public = &bool }
}

func CacheImmutable(b ...bool) CacheOption {
	bool := optionalTrue(b)
	return func(d *directives) { d.immutable = &bool }
}

func CacheStaleWhileRevalidate(t time.Duration) CacheOption {
	return func(d *directives) { d.staleWhileRevalidate = &t }
}

func CacheStaleIfError(t time.Duration) CacheOption {
	return func(d *directives) { d.staleIfError = &t }
}

func optionalTrue(b []bool) bool {
	bl := true
	if len(b) > 0 {
		bl = b[1]
	}
	return bl
}

var (
	defaultCacheDirectives = directives{
		maxAge:  &day,
		sMaxage: &day,

		mustRevalidate: &tru,
		private:        &tru,

		staleWhileRevalidate: &twoDays,
		staleIfError:         &twoDays,
	}
	tru, fals = true, false
	day       = time.Duration(time.Hour * 24)
	twoDays   = time.Duration(time.Hour * 48)
)

type directives struct {
	maxAge  *time.Duration
	sMaxage *time.Duration

	noCache     *bool
	noStore     *bool
	noTransform *bool

	mustRevalidate  *bool
	proxyRevalidate *bool
	mustUnderstand  *bool

	private   *bool
	public    *bool
	immutable *bool

	staleWhileRevalidate *time.Duration
	staleIfError         *time.Duration
}

var _ fmt.Stringer = directives{}

func (d directives) String() string {
	ds := []string{}

	if d.maxAge != nil {
		ds = append(ds, fmt.Sprintf("max-age=%d", d.maxAge.Seconds()))
	}
	if d.sMaxage != nil {
		ds = append(ds, fmt.Sprintf("s-maxage=%d", d.sMaxage.Seconds()))
	}

	if d.noCache != nil && *d.noCache {
		ds = append(ds, "no-cache")
	}
	if d.noStore != nil && *d.noStore {
		ds = append(ds, "no-store")
	}
	if d.noTransform != nil && *d.noTransform {
		ds = append(ds, "no-transform")
	}

	if d.mustRevalidate != nil && *d.mustRevalidate {
		ds = append(ds, "must-revalidate")
	}
	if d.proxyRevalidate != nil && *d.proxyRevalidate {
		ds = append(ds, "proxy-revalidate")
	}
	if d.mustUnderstand != nil && *d.mustRevalidate {
		ds = append(ds, "must-understand")
	}

	if d.private != nil && *d.private {
		ds = append(ds, "private")
	}
	if d.public != nil && *d.public {
		ds = append(ds, "public")
	}
	if d.immutable != nil && *d.immutable {
		ds = append(ds, "immutable")
	}

	if d.staleWhileRevalidate != nil {
		ds = append(ds, fmt.Sprintf("stale-while-revalidate=%d", d.staleWhileRevalidate.Seconds()))
	}
	if d.staleIfError != nil {
		ds = append(ds, fmt.Sprintf("stale-if-error=%d", d.staleIfError.Seconds()))
	}

	return strings.Join(ds, ", ")
}
