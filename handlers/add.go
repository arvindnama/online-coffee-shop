package handlers

import (
	"build-go-microservice/data"
	"net/http"
)

// swagger:route POST /products products createProduct
// Add a product into the database
// consumes:
//
// responses:
//
//	200: noContentResponse
func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	prod := r.Context().Value(KeyProduct{}).(*data.Product) // this is how we cast

	data.AddProduct(prod)
}
