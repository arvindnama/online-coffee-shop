// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import (
	"build-go-microservice/data"
	"context"
	"fmt"
	"log"
	"net/http"
)

// A list of products
// swagger:response productsResponse
type ProductsResponse struct {
	// All products in the system
	// in: body
	Body []data.Products
}

// swagger:parameters deleteProduct
type productIDPathParameterWrapper struct {
	// The id of the product to delete
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

type Products struct {
	l *log.Logger
}

type KeyProduct struct {
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) MiddlewareValidateProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJson(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(rw, "Unable to deserialize Product json", http.StatusBadRequest)
			return
		}

		err = prod.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				rw,
				fmt.Sprintf("Error validating product: %s", err),
				http.StatusBadRequest,
			)
			return
		}

		// add product to the context
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		r = r.WithContext(ctx)

		// call the next handler
		next.ServeHTTP(rw, r)
	})
}
