package handlers

import (
	"build-go-microservice/data"
	"net/http"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//	200: productsResponse

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	products := data.GetProducts()
	err := data.ToJSON(products, rw)

	if err != nil {
		http.Error(rw, "Unable to marshal produces", http.StatusInternalServerError)
	}

}
