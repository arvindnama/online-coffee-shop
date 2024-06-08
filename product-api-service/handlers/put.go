package handlers

import (
	"net/http"
	"product-api-service/data"
)

// swagger:route PUT /products/{id} products updateProduct
// Update a products details
//
// responses:
//
//	201: NoContentResponse
//	404: ErrorResponse
//	422: ErrorValidation
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Update product")

	id := getProductID(r)

	// [learning]: 1. Gorilla mux middleware store Data in Context
	// We read the product from the context which is stored with key `KeyProduct`
	// [learning]: 2. <>.(<>) is the syntax for typecasting
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err := data.UpdateProduct(id, prod)

	if err == data.ErrPrdNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}
