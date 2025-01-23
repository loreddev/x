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

package core

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"

	"forge.capytal.company/loreddev/x/blogo/plugin"
	"forge.capytal.company/loreddev/x/tinyssert"
)

// Creates a implementation of [http.Handler] that maps the [(*http.Request).Path] to a file of the
// same name in the file system provided by the sourcer. Use [Opts] to have more fine grained control
// over some additional behaviour of the implementation.
func NewServer(sourcer plugin.Sourcer, renderer plugin.Renderer, opts ...Opts) http.Handler {
	opt := Opts{}
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.Assertions == nil {
		opt.Assertions = tinyssert.NewDisabledAssertions()
	}
	if opt.Logger == nil {
		opt.Logger = slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{}))
	}
	if opt.TemplateErr == nil {
		opt.TemplateErr = templateErr
	}

	var filesystem fs.FS
	if opt.SourceOnInit {
		fs, err := sourcer.Source()
		if err != nil {
			panic(fmt.Sprintf("Failed to source files on initialization due to error: %s",
				err.Error(),
			))
		}
		filesystem = fs
	}

	return &server{
		files:       filesystem,
		sourcer:     sourcer,
		renderer:    renderer,
		assert:      opt.Assertions,
		log:         opt.Logger,
		errTemplate: opt.TemplateErr,
	}
}

// Options used in the construction of the server/[http.Handler] in [NewServer] to better
// control additional behaviour of the implementation.
type Opts struct {
	// Call [(plugin.Sourcer).Source] on construction of the implementation on [NewServer]?
	// Panics if the it returns a error. By default sourcing of files is done on the first
	// request.
	SourceOnInit bool
	// [tinyssert.Assertions] implementation used by server for it's Assertions, by default
	// uses [tinyssert.NewDisabledAssertions] to effectively disable assertions. Use this
	// if you want to the server to fail-fast on incorrect states.
	Assertions tinyssert.Assertions
	// Logger to be used to send error, warns and debug messages, useful for plugin development
	// and debugging the pipeline of files. By default it uses a logger that writes to [io.Discard],
	// effectively disabling logging.
	Logger *slog.Logger
	// Template used when the handler needs to return a non-200 status code. It is executed with
	// [ServeError] as data. Uses by default a plain text template.
	TemplateErr *template.Template
}

var templateErr = template.Must(template.New("defaultTemplateErr").Parse(
	"{{.StatusCode}} {{.Path}} failed to serve file {{.FileName}} {{.ErrMessage}}\n{{.Err}}",
))

type server struct {
	files fs.FS

	sourcer  plugin.Sourcer
	renderer plugin.Renderer

	assert      tinyssert.Assertions
	log         *slog.Logger
	errTemplate *template.Template
}

func (srv *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	srv.assert.NotNil(srv.sourcer, "A sourcer needs to be available to serve files")
	srv.assert.NotNil(srv.renderer, "A renderer needs to be available to serve files")
	srv.assert.NotNil(srv.log)
	srv.assert.NotNil(w)
	srv.assert.NotNil(r)

	log := srv.log.With(slog.String("path", r.URL.Path))
	log.Debug("Serving endpoint")

	if srv.files == nil {
		err := srv.serveHTTPSource(w, r)
		if err != nil {
			return
		}
	}

	path := strings.Trim(r.URL.Path, "/")
	if path == "" || path == "/" {
		path = "."
	}

	file, err := srv.serveHTTPOpenFile(path, w, r)
	if err != nil {
		return
	}

	// Defers the closing of the file to prevent memory being held if a renderer
	// does not properly closes the file.
	defer file.Close()

	err = srv.serveHTTPRender(file, w, r)
	if err != nil {
		return
	}

	log.Debug("Finished serving endpoint")
}

func (srv *server) serveHTTPSource(w http.ResponseWriter, r *http.Request) error {
	srv.assert.NotNil(srv.sourcer, "A sourcer needs to be available")
	srv.assert.NotNil(srv.errTemplate, "An error template needs to be available in cases of errors")
	srv.assert.NotNil(srv.log)
	srv.assert.NotNil(w)
	srv.assert.NotNil(r)

	log := srv.log.With(slog.String("path", r.URL.Path), slog.String("sourcer", srv.sourcer.Name()))
	log.Debug("Initializing file system")

	fs, err := srv.sourcer.Source()
	if err != nil {
		log.Error(
			"Failed to get file system, returning 500 code",
			slog.String("err", err.Error()),
		)

		w.WriteHeader(http.StatusInternalServerError)

		if err := srv.errTemplate.Execute(w, &ServeError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			ErrMessage: err.Error(),
			Path:       r.URL.Path,
		}); err != nil {
			log.Error("Failed to use error template", slog.String("err", err.Error()))
			_, err = w.Write([]byte(err.Error()))
			srv.assert.Nil(err)
		}

		return err
	}

	srv.files = fs

	return nil
}

func (srv *server) serveHTTPOpenFile(
	name string,
	w http.ResponseWriter,
	r *http.Request,
) (fs.File, error) {
	srv.assert.NotZero(name, "Name of file should not be empty")
	srv.assert.NotNil(srv.files, "A file system needs to be present to open a file")
	srv.assert.NotNil(srv.errTemplate, "An error template needs to be available in cases of errors")
	srv.assert.NotNil(srv.log)
	srv.assert.NotNil(w)
	srv.assert.NotNil(r)

	log := srv.log.With(
		slog.String("path", r.URL.Path),
		slog.String("filename", name),
		slog.String("sourcer", srv.sourcer.Name()),
	)
	log.Debug("Opening file")

	f, err := srv.files.Open(name)

	if errors.Is(err, fs.ErrNotExist) {
		log.Warn("File does not exists, returning 404 code",
			slog.String("err", err.Error()),
		)

		w.WriteHeader(http.StatusNotFound)

		if err := srv.errTemplate.Execute(w, &ServeError{
			StatusCode: http.StatusNotFound,
			Err:        err,
			ErrMessage: err.Error(),
			Path:       r.URL.Path,
			FileName:   name,
		}); err != nil {
			_, err = w.Write([]byte(err.Error()))
			srv.assert.Nil(err)
		}

		return nil, err
	} else if err != nil {
		log.Error("Failed to open file, returning 500 code",
			slog.String("err", err.Error()),
		)

		w.WriteHeader(http.StatusInternalServerError)

		if err := srv.errTemplate.Execute(w, &ServeError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			ErrMessage: err.Error(),
			Path:       r.URL.Path,
			FileName:   name,
		}); err != nil {
			_, err = w.Write([]byte(err.Error()))
			srv.assert.Nil(err)
		}

		return nil, err
	} else if f == nil {
		log.Error("File system returned a nil file, returning 500 code")

		w.WriteHeader(http.StatusInternalServerError)

		err := fmt.Errorf("file system returned a nil file using sourcer %q", srv.sourcer.Name())
		if err := srv.errTemplate.Execute(w, &ServeError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			ErrMessage: err.Error(),
			Path:       r.URL.Path,
			FileName:   name,
		}); err != nil {
			_, err = w.Write([]byte(err.Error()))
			srv.assert.Nil(err)
		}

		return nil, err
	}

	return f, err
}

func (srv *server) serveHTTPRender(file fs.File, w http.ResponseWriter, r *http.Request) error {
	srv.assert.NotNil(file, "A file needs to be present to it to be rendered")
	srv.assert.NotNil(srv.renderer, "A renderer needs to be present to render a file")
	srv.assert.NotNil(srv.errTemplate, "An error template needs to be available in cases of errors")
	srv.assert.NotNil(srv.log)
	srv.assert.NotNil(w)
	srv.assert.NotNil(r)

	log := srv.log.With(
		slog.String("path", r.URL.Path),
		slog.String("renderer", srv.renderer.Name()),
	)
	log.Debug("Rendering file")

	err := srv.renderer.Render(file, w)
	if err != nil {
		log.Error("Failed to render file, returning 500 code")

		w.WriteHeader(http.StatusInternalServerError)

		if err := srv.errTemplate.Execute(w, &ServeError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
			ErrMessage: err.Error(),
			Path:       r.URL.Path,
		}); err != nil {
			_, err = w.Write([]byte(err.Error()))
			srv.assert.Nil(err)
		}

		return err
	}

	return nil
}

type ServeError struct {
	StatusCode int
	Err        error
	ErrMessage string
	Path       string
	FileName   string
}

func (e *ServeError) Error() string {
	return fmt.Sprintf("failed to serve file %q to endpoint %q", e.FileName, e.Path)
}
