package gohserver

import "net/http"

// StatusWriter defines a http response writer for prometheus
type StatusWriter struct {
	http.ResponseWriter
	StatusCode int
}

// Init initializes the StatusWriter
func (w *StatusWriter) Init(o http.ResponseWriter) *StatusWriter {
	w.ResponseWriter = o
	w.StatusCode = 200
	return w
}

// WriteHeader redefines the method for the StatusWriter wrapper
func (w *StatusWriter) WriteHeader(status int) {
	w.StatusCode = status
	w.ResponseWriter.WriteHeader(status)
}

// Write redefines the method for the StatusWriter wrapper
func (w *StatusWriter) Write(b []byte) (int, error) {
	return w.ResponseWriter.Write(b)
}

// Header redefines the methdo for the StatusWriter wrapper
func (w *StatusWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}
