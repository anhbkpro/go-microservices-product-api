package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

type GzipHandler struct {
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			wrw := NewWrappedResponseWriter(rw)
			wrw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(wrw, r)
			defer wrw.Flush()

			return
		}

		// handle normal
		next.ServeHTTP(rw, r)
	})
}

// WrappedResponseWriter creates custom ResponseWriter
type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(rw)
	return &WrappedResponseWriter{
		rw: rw,
		gw: gw,
	}
}

func (w *WrappedResponseWriter) Header() http.Header {
	return w.rw.Header()
}

func (w *WrappedResponseWriter) Write(d []byte) (int, error) {
	return w.gw.Write(d)
}

func (w *WrappedResponseWriter) WriteHeader(statusCode int) {
	w.rw.WriteHeader(statusCode)
}

func (w *WrappedResponseWriter) Flush() {
	w.gw.Flush()
	w.gw.Close()
}
