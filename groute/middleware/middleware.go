package middleware

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Middleware = func(next http.Handler) http.Handler

type MiddlewaredReponse struct {
	w          http.ResponseWriter
	statuses   []int
	bodyWrites [][]byte
}

func NewMiddlewaredResponse(w http.ResponseWriter) *MiddlewaredReponse {
	return &MiddlewaredReponse{w, []int{500}, [][]byte{[]byte("")}}
}

func (m *MiddlewaredReponse) WriteHeader(s int) {
	m.Header().Set("Status", strconv.Itoa(s))
	m.statuses = append(m.statuses, s)
}

func (m *MiddlewaredReponse) Header() http.Header {
	return m.w.Header()
}

func (m *MiddlewaredReponse) Write(b []byte) (int, error) {
	m.bodyWrites = append(m.bodyWrites, b)
	return len(b), nil
}

func (m *MiddlewaredReponse) ReallyWriteHeader() (int, error) {
	status := m.statuses[len(m.statuses)-1]
	m.w.WriteHeader(status)
	bytes := 0
	for _, b := range m.bodyWrites {
		by, err := m.w.Write(b)
		if err != nil {
			return bytes, errors.Join(
				fmt.Errorf(
					"Failed to write to response in middleware."+
						"\nStatuses are %v"+
						"\nTried to write %v bytes"+
						"\nTried to write response:\n%s",
					m.statuses, bytes, string(b),
				),
				err,
			)
		}
		bytes += by
	}

	return bytes, nil
}

type multiResponseWriter struct {
	response http.ResponseWriter
	writers  []io.Writer
}

func MultiResponseWriter(
	w http.ResponseWriter,
	writers ...io.Writer,
) http.ResponseWriter {
	if mw, ok := w.(*multiResponseWriter); ok {
		mw.writers = append(mw.writers, writers...)
		return mw
	}

	allWriters := make([]io.Writer, 0, len(writers))
	for _, iow := range writers {
		if mw, ok := iow.(*multiResponseWriter); ok {
			allWriters = append(allWriters, mw.writers...)
		} else {
			allWriters = append(allWriters, iow)
		}
	}

	return &multiResponseWriter{w, allWriters}
}

func (w *multiResponseWriter) WriteHeader(status int) {
	w.Header().Set("Status", strconv.Itoa(status))
	w.response.WriteHeader(status)
}

func (w *multiResponseWriter) Write(p []byte) (int, error) {
	w.WriteHeader(http.StatusOK)
	for _, w := range w.writers {
		n, err := w.Write(p)
		if err != nil {
			return n, err
		}
		if n != len(p) {
			return n, io.ErrShortWrite
		}
	}
	return w.response.Write(p)
}

func (w *multiResponseWriter) Header() http.Header {
	return w.response.Header()
}
