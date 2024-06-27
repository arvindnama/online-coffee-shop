package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type WrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *WrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// [learning]: to set the base class, we can set by using type of the base class.
		// in TS we use spread operator ({..obj}), in go it is ({BaseType : obj})
		ww := &WrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusAccepted,
		}
		next.ServeHTTP(ww, r)

		m.logger.Info(fmt.Sprintf("%d", ww.statusCode), r.URL.Path, time.Since(start))

	})
}
