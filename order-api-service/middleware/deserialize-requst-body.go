package middleware

import (
	"context"
	"net/http"

	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/arvindnama/golang-microservices/order-service/data"
)

type RequestBody struct{}

func (m *Middleware) DeserializeRequestBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Trace("Deserializing request body")
		defer r.Body.Close()
		var order data.Order
		dataUtils.FromJSON(&order, r.Body)

		ctx := context.WithValue(r.Context(), RequestBody{}, order)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
