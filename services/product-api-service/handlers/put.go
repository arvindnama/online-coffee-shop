package handlers

import (
	"net/http"

	"github.com/arvindnama/golang-microservices/product-api-service/data"
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
	p.l.Info("Handle Update product")

	id := getProductID(r)

	// [learning]: 1. Gorilla mux middleware store Data in Context
	// We read the product from the context which is stored with key `KeyProduct`
	// [learning]: 2. <>.(<>) is the syntax for typecasting
	prod := r.Context().Value(KeyProduct{}).(*data.Product)

	err := p.pDB.UpdateProduct(id, prod)

	if err == data.ErrPrdNotFound {
		p.l.Error("Product not found", err)
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		p.l.Error("Product not found", err)
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}
