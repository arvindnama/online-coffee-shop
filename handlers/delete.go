package handlers

import (
	"build-go-microservice/data"
	"net/http"
)

// swagger:route DELETE /products/{id} products deleteProduct
// Deletes the product from the database
// responses:
//
//	200: noContentResponse
//	404: errorResponse
//	501: errorResponse
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
