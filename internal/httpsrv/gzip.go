package httpsrv

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	gz         io.Writer
	origWriter http.ResponseWriter
}

func (w gzipResponseWriter) Write(b []byte) (int, error) {
	return w.gz.Write(b)
}

func (w gzipResponseWriter) WriteHeader(statusCode int) {
	w.origWriter.WriteHeader(statusCode)
}

func (w gzipResponseWriter) Header() http.Header {
	return w.origWriter.Header()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Add("Vary", "Accept-Encoding")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		gzw := gzipResponseWriter{gz: gz, origWriter: w}
		next.ServeHTTP(gzw, r)
	})
}
