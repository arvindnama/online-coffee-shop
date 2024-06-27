package middleware

import (
	"fmt"
	"net/http"

	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/arvindnama/golang-microservices/order-service/data"
)

func (m *Middleware) ValidateRequestBody(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.Trace("Validating Request body")
		order := r.Context().Value(RequestBody{}).(data.Order)

		m.logger.Debug(fmt.Sprintf("%#v", order))
		errs := m.validation.Validate(order)
		if errs != nil {
			m.logger.Error("[Error] validating order", errs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			dataUtils.ToJSON(
				&data.ValidationError{
					Messages: errs.Errors(),
				}, w,
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}
