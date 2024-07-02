package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type GzipHandler struct {
	logger hclog.Logger
}

type WrappedResponseWriter struct {
	rw http.ResponseWriter
	gw *gzip.Writer
}

func NewGzipHandler(logger hclog.Logger) *GzipHandler {
	return &GzipHandler{logger}
}

func NewWrappedResponseWrite(rw http.ResponseWriter) *WrappedResponseWriter {
	gw := gzip.NewWriter(rw)
	return &WrappedResponseWriter{rw, gw}
}

func (wr *WrappedResponseWriter) Header() http.Header {
	return wr.rw.Header()
}

func (wr *WrappedResponseWriter) Write(data []byte) (int, error) {
	return wr.gw.Write(data)
}

func (wr *WrappedResponseWriter) WriteHeader(statusCode int) {
	wr.rw.WriteHeader(statusCode)
}

func (wr *WrappedResponseWriter) Flush() {
	wr.gw.Flush()
}

func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			g.logger.Debug("Requests accepts Gzip encoding, Response will be giz compressed")

			// create gziped response
			wwr := NewWrappedResponseWrite(rw)
			wwr.Header().Add("Content-Encoding", "gzip")
			next.ServeHTTP(wwr, r)
			//[learning]: defer is a key in go , to tell go execute the statement
			// at the end of function , something like finally{}
			defer wwr.Flush()
			return
		}
		// call the next handler
		next.ServeHTTP(rw, r)
	})
}
