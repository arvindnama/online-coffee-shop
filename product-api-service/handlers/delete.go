package handlers

import (
	"net/http"

	"github.com/arvindnama/golang-microservices/product-api-service/data"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes the product from the database
// responses:
//
//	200: NoContentResponse
//	404: ErrorResponse
//	501: ErrorResponse
func (p *Products) DeleteProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Delete product")

	id := getProductID(r)

	err := data.DeleteProduct(id)

	if err == data.ErrPrdNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}
