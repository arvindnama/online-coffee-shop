package middleware

import (
	"net/http"

	"github.com/gorilla/handlers"
)

func (m *Middleware) AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsHandler := handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(next)

		m.logger.Info("CORS Handler called")
		corsHandler.ServeHTTP(w, r)
	})
}
