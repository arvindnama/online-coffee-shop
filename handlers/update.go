package handlers

import (
	"build-go-microservice/data"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// swagger:route PUT /products products updateProduct
// Update a products details
//
// responses:
//
//	201: noContentResponse
func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Update product")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
		return
	}

	prod := r.Context().Value(KeyProduct{}).(*data.Product) // this is how we cast

	err = data.UpdateProduct(id, prod)

	if err == data.ErrPrdNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}

}
