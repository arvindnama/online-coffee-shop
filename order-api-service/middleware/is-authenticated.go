package middleware

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"strings"
)

type AuthUserID struct{}

func (m *Middleware) writeUnauthenticated(err error, w http.ResponseWriter) {
	m.logger.Error(err.Error())
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}
func (m *Middleware) IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer") {
			err := errors.New("authorization token not present")
			m.writeUnauthenticated(err, w)
			return

		}
		encodedToken := strings.TrimPrefix(authorization, "Bearer ")
		token, err := base64.StdEncoding.DecodeString(encodedToken)
		if err != nil {
			m.writeUnauthenticated(err, w)
			return
		}
		m.logger.Info(r.URL.Path, "Authorized")
		userId := string(token)

		ctx := context.WithValue(r.Context(), AuthUserID{}, userId)
		neReq := r.WithContext(ctx)

		next.ServeHTTP(w, neReq)
	})
}
