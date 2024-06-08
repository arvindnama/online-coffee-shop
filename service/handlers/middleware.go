package handlers

import (
	"build-go-microservice/data"
	"context"
	"net/http"
)

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := data.FromJSON(prod, r.Body)
		p.l.Println(prod)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to deserialize Product json", http.StatusBadRequest)
			return
		}

		errs := p.v.Validate(prod)
		if errs != nil {
			p.l.Println("[ERROR] validating product", err)

			//return errors array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			data.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		// add product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(rw, r)
	})
}
