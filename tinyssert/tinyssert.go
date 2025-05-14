// Copyright (c) 2025 Gustavo "Guz" L. de Mello
// Copyright (c) 2025 The Lored.dev Contributors
//
// Contents of this file, expect as otherwise noted, are dual-licensed under the
// Apache License, Version 2.0 <http://www.apache.org/licenses/LICENSE-2.0> or
// the MIT license <http://opensource.org/licenses/MIT>, at you option.
//
// You may use this file in compliance with the licenses.
//
// Unless required by applicable law or agreed to in writing, this file distributed
// under the licenses is distributed on as "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS
// OF ANY KIND, either express or implied.
//
// An original copy of this file can be found at http://forge.capytal.company/loreddev/x/tinyssert/tinyssert.go.

// Package tinyssert is a minimal set of assertions functions for testing and simulation
// testing, all in one file.
//
// The most simple way of using the package is importing it directly and using the
// alias functions:
//
//	package main
//
//	import (
//	  "log"
//	  "forge.capytal.company/loreddev/x/tinyssert"
//	)
//
//	func main() {
//	  expected := "value"
//	  value := "value"
//	  log.Println(tinyssert.Equal(expected, value)) // "true"
//	}
//
// You can create your own "assert" variable and have more control
// over how asserts work, see the [New] constructor for more information:
//
//	package main
//
//	import (
//	  "log/slog"
//	  "forge.capytal.company/loreddev/x/tinyssert"
//	)
//
//	var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
//	var assert = tinyssert.NewAssertions(tinyssert.WithLogger(logger))
//
//	func main() {
//	  expected := "value"
//	  value := "not value"
//	  assert.Equal(expected, value) // "expected \"value\", got \"not value\"" with the call stack and returns false
//	}
//
// Preferably, when using assertions inside production code or libraries, you can use
// the assertions via dependency injection. This provides a easy way to disable
// assertions in production (see [NewDisabled]) while being able to test an API without
// changing it:
//
//	package main
//
//	import (
//	  "flag"
//	  "log/slog"
//
//	  "forge.capytal.company/loreddev/x/tinyssert"
//	)
//
//	var debug = flag.Bool("debug", false, "Run the application in debug mode")
//
//	func init() {
//	  flag.Parse()
//	}
//
//	func main() {
//	  logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
//	  assert := tinyssert.NewDisabled()
//	  if *debug {
//	    assert := tinyssert.New(tinyssert.WithLogger(logger))
//	  }
//
//	  app := App{logger: logger, assert: assert}
//
//	  app.Start()
//	}
//
//	type App struct {
//	  logger: *slog.Logger
//	  assert: tinyssert.Assertions
//	}
//
//	function (app *App) Start() {
//	  app.assert.OK(app.logger, "Logger must be initialized before the application")
//
//	  // ...
//	}
package tinyssert

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Assertions represents all the API provided by the [assert] package. Implementation
// of this interface can have extra logic on their state such as panicking or using
// [testing.T.FailNow] to halt the program from continuing on failed assertions instead
// of just return false.
//
// `msg` argument(s) is(/are) optional, the first argument should always be a string.
// The string can have formatting verbs, with the remaining arguments being used to fill
// those verbs. So for example:
//
//	 if argument > 0 {
//		  tinyssert.OK(value, "Since %d is greater than 0, this should always be ok", argument)
//	 }
type Assertions interface {
	// Asserts that the value is not zero-valued, is nil, or panics, aka. "is ok".
	OK(v any, msg ...any) error

	// Asserts that the actual value is equal to the expected value.
	Equal(expected, actual any, msg ...any) error
	// Asserts that the actual value is not equal to the expected value.
	NotEqual(notExpected, actual any, msg ...any) error

	// Asserts that the value is nil.
	Nil(v any, msg ...any) error
	// Asserts that the value is not nil.
	NotNil(v any, msg ...any) error

	// Asserts that the value is a boolean true.
	True(b bool, msg ...any) error
	// Asserts that the value is a boolean false.
	False(b bool, msg ...any) error

	// Asserts that the value is zero-valued.
	Zero(v any, msg ...any) error
	// Asserts that the value is not zero-valued.
	NotZero(v any, msg ...any) error

	// Asserts that the function panics.
	Panic(fn func(), msg ...any) error
	// Asserts that the function does not panics.
	NotPanic(fn func(), msg ...any) error

	// Logs the formatted failure message and/or marks the test as failed if possible,
	// depending of what is possible to the implementation.
	Fail(f Failure)
	// Panics with the formatted failure message and/or marks the test as failed,
	// depending of what is possible to the implementation.
	FailNow(f Failure)

	// Gets the caller stack.
	CallerInfo() []string
}

// New constructs a new implementation of [Assertions]. Use `opts` to customize the behaviour
// of the implementation.
func New(opts ...Option) Assertions {
	a := &assertions{
		panic: false,
		log:   slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})),

		test:   nil,
		helper: nil,
	}

	for _, opt := range opts {
		opt(a)
	}

	if th, ok := a.test.(helperT); ok {
		a.helper = th
	}

	return a
}

type Option = func(*assertions)

func WithPanic() Option {
	return func(a *assertions) {
		a.panic = true
	}
}

func WithTest(t TestingT) Option {
	return func(a *assertions) {
		a.test = t
	}
}

func WithLogger(l *slog.Logger) Option {
	return func(a *assertions) {
		a.log = l
	}
}

type assertions struct {
	panic bool

	test   TestingT
	helper helperT

	log   *slog.Logger
	group string
}

// Wrapper interface around [testing.T].
type TestingT interface {
	Errorf(format string, args ...any)
}
type helperT interface {
	Helper()
}

var _ Assertions = (*assertions)(nil)

func (a *assertions) Equal(expected, actual any, msg ...any) error {
	if a.equal(expected, actual) {
		return nil
	}
	return a.fail(fmt.Sprintf("expected %v (right), got %v (left)", expected, actual), msg...)
}

func (a *assertions) NotEqual(notExpected, actual any, msg ...any) error {
	if !a.equal(notExpected, actual) {
		return nil
	}
	return a.fail(fmt.Sprintf("expected to %v (right) and %v (left) to be not-equal", notExpected, actual), msg...)
}

func (a *assertions) equal(ex, ac any) bool {
	if nex, nac := a.OK(ex), a.OK(ac); (nex != nil) != (nac != nil) {
		return false
	}

	if reflect.DeepEqual(ex, ac) {
		return true
	}

	ev, av := reflect.ValueOf(ex), reflect.ValueOf(ac)

	if ev == av {
		return true
	}

	if av.Type().ConvertibleTo(ev.Type()) {
		return reflect.DeepEqual(ex, av.Convert(ev.Type()).Interface())
	}

	if fmt.Sprintf("%#v", ex) == fmt.Sprintf("%#v", ac) {
		return true
	}

	return false
}

func (a *assertions) OK(v any, msg ...any) error {
	if a.nil(v) {
		return a.fail("expected not-nil value", msg...)
	}
	if a.zero(v) {
		return a.fail("expected non-zero value", msg...)
	}

	if f, ok := v.(func()); ok {
		if a.panics(f) {
			return a.fail("expected to not panic")
		}
	}

	return nil
}

func (a *assertions) Nil(v any, msg ...any) error {
	if a.nil(v) {
		return nil
	}
	return a.fail("expected nil value", msg...)
}

func (a *assertions) NotNil(v any, msg ...any) error {
	if !a.nil(v) {
		return nil
	}
	return a.fail("expected not-nil value", msg...)
}

func (a *assertions) nil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	rk := rv.Kind()
	if rk >= reflect.Chan && rk <= reflect.Slice && rv.IsNil() {
		return true
	}

	return false
}

func (a *assertions) True(v bool, msg ...any) error {
	if v {
		return nil
	}
	return a.fail("expected true", msg...)
}

func (a *assertions) False(v bool, msg ...any) error {
	if !v {
		return nil
	}
	return a.fail("expected false", msg...)
}

func (a *assertions) Zero(v any, msg ...any) error {
	if a.zero(v) {
		return nil
	}
	return a.fail("expected zero value", msg...)
}

func (a *assertions) NotZero(v any, msg ...any) error {
	if !a.zero(v) {
		return nil
	}
	return a.fail("expected non-zero value", msg...)
}

func (a *assertions) zero(v any) bool {
	if v != nil && !reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
		return false
	}
	return true
}

func (a *assertions) Panic(fn func(), msg ...any) error {
	if a.panics(fn) {
		return nil
	}
	return a.fail("expected function to panic", msg...)
}

func (a *assertions) NotPanic(fn func(), msg ...any) error {
	if !a.panics(fn) {
		return nil
	}
	return a.fail("expected function to not panic", msg...)
}

func (a *assertions) panics(fn func()) bool {
	var r any
	func() {
		defer func() {
			r = recover()
		}()
		fn()
	}()
	return r != nil
}

func (a *assertions) fail(reason string, msg ...any) error {
	if a.helper != nil {
		a.helper.Helper()
	}

	f := Failure{
		Reason:     reason,
		Message:    fmtMessage(msg),
		CallerInfo: a.CallerInfo(),
	}

	if n, ok := a.test.(interface {
		Name() string
	}); ok {
		f.Test = n.Name()
	}

	if a.panic {
		a.FailNow(f)
	} else {
		a.Fail(f)
	}

	return f
}

func (a *assertions) Fail(f Failure) {
	if ft, ok := a.test.(interface {
		Fail()
	}); ok {
		a.test.Errorf("ASSERTION FAILED:\n%s", f.String())
		ft.Fail()
	} else {
		a.log.Error("ASSERTION FAILED",
			slog.String("reason", f.Reason),
			slog.String("message", f.Message),
			slog.String("test", f.Test),
			slog.Any("caller", f.CallerInfo),
		)
	}
}

func (a *assertions) FailNow(f Failure) {
	if ft, ok := a.test.(interface {
		FailNow()
	}); ok {
		a.test.Errorf("ASSERTION FAILED:\n%s", f.String())
		ft.FailNow()
	} else {
		panic(f.String())
	}
}

func fmtMessage(msg ...any) string {
	switch len(msg) {
	case 0:
		return ""
	case 1:
		if s, ok := msg[0].(string); ok {
			return s
		}
		return fmt.Sprintf("%v", msg[0])
	default:
		var m string
		if s, ok := msg[0].(string); ok {
			m = s
		} else {
			m = fmt.Sprintf("%v", msg[0])
		}
		return fmt.Sprintf(m, msg[1:]...)
	}
}

func (as *assertions) CallerInfo() []string {
	callers := []string{}
	for i := 0; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			// We reached the end of the call stack
			break
		}

		// Edge case found in https://github.com/stretchr/testify/issues/180
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}
		name := f.Name()

		if name == "testing.Runner" {
			break
		}

		filename := path.Base(file)
		dirname := path.Base(path.Dir(file))
		if (dirname != "assert" && dirname != "mock" && dirname != "require") ||
			filename == "mock_test.go" {
			callers = append(callers, fmt.Sprintf("%s:%d", file, line))
		}

		// Remove the package
		s := strings.Split(name, ".")
		name = s[len(s)-1]

		if isTest(name, "Test") || isTest(name, "Benchmark") || isTest(name, "Example") {
			break
		}
	}
	return callers
}

func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}
	if len(name) == len(prefix) {
		return true
	}
	r, _ := utf8.DecodeRuneInString(name[len(prefix):])
	return !unicode.IsLower(r)
}

type Failure struct {
	Reason  string
	Message string

	Test       string
	CallerInfo []string
}

var (
	_ error        = Failure{}
	_ fmt.Stringer = Failure{}
)

func (e Failure) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("assertion failed, %s: %s", e.Reason, e.Message)
	}
	return fmt.Sprintf("assertion failed, %s", e.Reason)
}

func (e Failure) String() string {
	c := map[string]string{
		"Reason": e.Reason,
	}

	if e.Message != "" {
		c["Message"] = e.Message
	}

	if e.Test != "" {
		c["Test"] = e.Test
	}

	c["Stack Trace"] = e.StackTrace()

	var out string
	for k, m := range c {
		var s string
		for _, l := range strings.Split(m, "\n") {
			s += fmt.Sprintf("\t%s\n", l)
		}
		out += fmt.Sprintf("\t%s:\n%s", k, s)
	}

	return out
}

// StackTrace returns the CallerInfo strings as a formatted stack trace.
func (e Failure) StackTrace() string {
	return strings.Join(e.CallerInfo, "\n\t")
}

type disabledAssertions struct{}

// NewDisabled creates a new implementation of Assertions that always a nil error and
// never panics or marks the test as failed, with the exception of Fail, FailNow and
// CallerInfo, which uses their corresponding [Fail], [FailNow] and [CallerInfo] top-level
// functions.
//
// The `opts` argument does nothing, and is just available to make the function signature
// equal to [New].
func NewDisabled(opts ...Option) Assertions {
	_ = opts
	return &disabledAssertions{}
}

func (*disabledAssertions) OK(any, ...any) error              { return nil }
func (*disabledAssertions) Equal(_, _ any, _ ...any) error    { return nil }
func (*disabledAssertions) NotEqual(_, _ any, _ ...any) error { return nil }
func (*disabledAssertions) Nil(any, ...any) error             { return nil }
func (*disabledAssertions) NotNil(any, ...any) error          { return nil }
func (*disabledAssertions) True(bool, ...any) error           { return nil }
func (*disabledAssertions) False(bool, ...any) error          { return nil }
func (*disabledAssertions) Zero(any, ...any) error            { return nil }
func (*disabledAssertions) NotZero(any, ...any) error         { return nil }
func (*disabledAssertions) Panic(func(), ...any) error        { return nil }
func (*disabledAssertions) NotPanic(func(), ...any) error     { return nil }
func (*disabledAssertions) Fail(f Failure)                    { Default.Fail(f) }
func (*disabledAssertions) FailNow(f Failure)                 { Default.FailNow(f) }
func (*disabledAssertions) CallerInfo() []string              { return Default.CallerInfo() }

var (
	// DefaultLogger is the default [slog.Logger] used by [Default]
	DefaultLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	// Default implementation of [Assertions] used by the top-level API functions.
	Default = New(WithLogger(DefaultLogger))
)

// OK asserts that the value is not zero-valued, is nil, or panics, aka. "is ok".
//
// Logs the failure message with [DefaultLogger].
func OK(obj any, msg ...any) error {
	return Default.OK(obj, msg...)
}

// Equal asserts that the actual value is equal to the expected value.
//
// Logs the failure message with [DefaultLogger].
func Equal(expected, actual any, msg ...any) error {
	return Default.Equal(expected, actual, msg...)
}

// NotEqual asserts that the actual value is not equal to the expected value.
//
// Logs the failure message with [DefaultLogger].
func NotEqual(notExpected, actual any, msg ...any) error {
	return Default.NotEqual(notExpected, actual, msg...)
}

// Nil asserts that the value is nil.
//
// Logs the failure message with [DefaultLogger].
func Nil(v any, msg ...any) error {
	return Default.Nil(v, msg...)
}

// NotNil asserts that the value is not nil.
//
// Logs the failure message with [DefaultLogger].
func NotNil(v any, msg ...any) error {
	return Default.NotNil(v, msg...)
}

// True asserts that the value is a boolean true.
//
// Logs the failure message with [DefaultLogger].
func True(v bool, msg ...any) error {
	return Default.True(v, msg...)
}

// False asserts that the value is a boolean false.
//
// Logs the failure message with [DefaultLogger].
func False(v bool, msg ...any) error {
	return Default.False(v, msg...)
}

// Zero asserts that the value is zero-valued.
//
// Logs the failure message with [DefaultLogger].
func Zero(v any, msg ...any) error {
	return Default.Zero(v, msg...)
}

// NotZero asserts that the value is not zero-valued.
//
// Logs the failure message with [DefaultLogger].
func NotZero(v any, msg ...any) error {
	return Default.NotZero(v, msg...)
}

// Panic asserts that the function panics.
//
// Logs the failure message with [DefaultLogger].
func Panic(fn func(), msg ...any) error {
	return Default.Panic(fn, msg...)
}

// NotPanic asserts that the function does not panics.
//
// Logs the failure message with [DefaultLogger].
func NotPanic(fn func(), msg ...any) error {
	return Default.NotPanic(fn, msg...)
}

// Fail logs the formatted failure message using [DefaultLogger].
func Fail(f Failure) {
	Default.Fail(f)
}

// FailNow panics with the formatted failure message.
func FailNow(f Failure) {
	Default.FailNow(f)
}
