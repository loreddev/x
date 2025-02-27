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
	"log/slog"
	"math/rand"
	"net"
	"net/http"
)

const (
	defaultMsgNew = "NEW REQUEST"
	defaultMsg200 = "END REQUEST"
	defaultMsg400 = "INV REQUEST"
	defaultMsg500 = "ERR REQUEST"
)

func Logger(logger *slog.Logger, options ...LoggerOption) Middleware {
	state := &loggerState{
		levelNew: slog.LevelDebug,
		msgNew:   "",
		argsNew:  LoggerArgsDefault,

		level200: -1,
		msg200:   "",
		args200:  LoggerArgsDefault,

		level400: -1,
		msg400:   "",
		args400:  LoggerArgsDefault,

		level500: -1,
		msg500:   "",
		args500:  LoggerArgsDefault,

		hashFunction: randHash,

		logger: logger,
	}

	for _, option := range options {
		option(state)
	}

	if state.level200 == -1 {
		state.level200 = state.levelNew + 4
	}
	if state.level400 == -1 {
		state.level400 = state.level200 + 4
	}
	if state.level500 == -1 {
		state.level500 = state.level500 + 4
	}

	if state.msgNew == "" {
		state.msgNew = defaultMsgNew
	}
	if state.msg200 == "" {
		if state.msgNew != "" && state.msgNew != defaultMsgNew {
			state.msg200 = state.msgNew
		} else {
			state.msg200 = defaultMsg200
		}
	}
	if state.msg400 == "" {
		if state.msg200 != "" && state.msg200 != defaultMsg200 {
			state.msg400 = state.msg200
		} else {
			state.msg400 = defaultMsg400
		}
	}
	if state.msg500 == "" {
		if state.msg400 != "" && state.msg400 != defaultMsg400 {
			state.msg500 = state.msg400
		} else {
			state.msg500 = defaultMsg500
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := randHash(5)

			state.logger = logger.With(slog.String("id", id))

			lw := &loggerResponseWriter{w, 0}

			state.lNew(state.argsNew(lw, r)...)

			next.ServeHTTP(lw, r)

			switch {
			case lw.StatusCode() >= 500:
				state.l500(state.args500(lw, r)...)
			case lw.StatusCode() >= 400:
				state.l400(state.args400(lw, r)...)
			default:
				state.l200(state.args200(lw, r)...)
			}
		})
	}
}

type LoggerOption func(*loggerState)

func LoggerWithLevel(l slog.Level) LoggerOption {
	return func(ls *loggerState) { ls.levelNew = l }
}

func LoggerWithMsg(msg string) LoggerOption {
	return func(ls *loggerState) { ls.msgNew = msg }
}

func LoggerWithArgs(args LoggerArgs) LoggerOption {
	return func(ls *loggerState) { ls.argsNew = args }
}

func LoggerWith200Level(l slog.Level) LoggerOption {
	return func(ls *loggerState) { ls.level200 = l }
}

func LoggerWith200Msg(msg string) LoggerOption {
	return func(ls *loggerState) { ls.msg200 = msg }
}

func LoggerWith200Args(args LoggerArgs) LoggerOption {
	return func(ls *loggerState) { ls.args200 = args }
}

func LoggerWith400Level(l slog.Level) LoggerOption {
	return func(ls *loggerState) { ls.level400 = l }
}

func LoggerWith400Msg(msg string) LoggerOption {
	return func(ls *loggerState) { ls.msg400 = msg }
}

func LoggerWith400Args(args LoggerArgs) LoggerOption {
	return func(ls *loggerState) { ls.args400 = args }
}

func LoggerWith500Level(l slog.Level) LoggerOption {
	return func(ls *loggerState) { ls.level500 = l }
}

func LoggerWith500Msg(msg string) LoggerOption {
	return func(ls *loggerState) { ls.msg500 = msg }
}

func LoggerWith500Args(args LoggerArgs) LoggerOption {
	return func(ls *loggerState) { ls.args500 = args }
}

type LoggerArgs func(LoggerResponseWriter, *http.Request) []any

func LoggerArgsDefault(lw LoggerResponseWriter, r *http.Request) []any {
	addr := LoggerGetAddr(r)

	if net.ParseIP(addr) == nil {
		addr = fmt.Sprintf("INVALID ADDR %s", addr)
	}

	return []any{
		slog.String("status", fmt.Sprintf("%3d", lw.StatusCode())),
		slog.String("method", fmt.Sprintf("%3s", r.Method)),
		slog.String("addr", addr),
		slog.String("path", r.URL.Path),
	}
}

func LoggerGetAddr(r *http.Request) string {
	if i := r.Header.Get("CF-Connecting-IP"); i != "" {
		return i
	}
	if i := r.Header.Get("X-Forwarded-For"); i != "" {
		return i
	}
	if i := r.Header.Get("X-Real-IP"); i != "" {
		return i
	}
	return r.RemoteAddr
}

type LoggerResponseWriter interface {
	http.ResponseWriter
	StatusCode() int
}

type loggerResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *loggerResponseWriter) WriteHeader(s int) {
	w.statusCode = s
	w.ResponseWriter.WriteHeader(s)
}

func (w *loggerResponseWriter) StatusCode() int {
	return w.statusCode
}

type loggerState struct {
	levelNew slog.Level
	msgNew   string
	argsNew  LoggerArgs

	level200 slog.Level
	msg200   string
	args200  LoggerArgs

	level400 slog.Level
	msg400   string
	args400  LoggerArgs

	level500 slog.Level
	msg500   string
	args500  LoggerArgs

	hashFunction func(n int) string

	logger *slog.Logger
}

func (l *loggerState) lNew(args ...any) {
	l.logLevel(l.levelNew, l.msgNew, args...)
}

func (l *loggerState) l200(args ...any) {
	l.logLevel(l.level200, l.msg200, args...)
}

func (l *loggerState) l400(args ...any) {
	l.logLevel(l.level400, l.msg400, args...)
}

func (l *loggerState) l500(args ...any) {
	l.logLevel(l.level500, l.msg500, args...)
}

func (l *loggerState) logLevel(level slog.Level, msg string, args ...any) {
	switch true {
	case level >= slog.LevelError:
		l.logger.Error(msg, args...)
	case level >= slog.LevelWarn:
		l.logger.Warn(msg, args...)
	case level >= slog.LevelInfo:
		l.logger.Info(msg, args...)
	default:
		l.logger.Debug(msg, args...)
	}
}

func getBiggestLength(s ...string) int {
	var l int
	for _, s := range s {
		if len(s) > l {
			l = len(s)
		}
	}
	return l
}

const HASH_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// This is not the most performant function, as a TODO we could
// improve based on this Stackoberflow thread:
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func randHash(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = HASH_CHARS[rand.Int63()%int64(len(HASH_CHARS))]
	}
	return string(b)
}
