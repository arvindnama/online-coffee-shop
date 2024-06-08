package handlers

import (
	"net/http"
	"product-api-service/data"
)

// swagger:route POST /products products createProduct
// Add a product into the database
// consumes:
//
// responses:
//
//	200: NoContentResponse
//	422: ErrorValidation
//	501: ErrorResponse
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	prod := r.Context().Value(KeyProduct{}).(*data.Product) // this is how we cast

	data.AddProduct(prod)
}
