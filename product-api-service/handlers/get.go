package handlers

import (
	"net/http"

	"github.com/arvindnama/golang-microservices/product-api-service/data"
)

// swagger:route GET /products products listProducts
// Returns a list of products
// responses:
//	200: ProductsResponse
//	501: ErrorResponse

func (p *Products) GetAllProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	products := data.GetProducts()
	rw.Header().Add("Content-Type", "application/json")
	err := data.ToJSON(products, rw)
	if err != nil {
		http.Error(rw, "Unable to marshal produces", http.StatusInternalServerError)
	}

}

// swagger:route GET /products/{id} products listProduct
// Returns a product
// responses:
//
//	200: ProductResponse
//	404: ErrorResponse
//	501: ErrorResponse
func (p *Products) GetProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	rw.Header().Add("Content-Type", "application/json")

	id := getProductID(r)
	product, err := data.GetProductById(id)

	if err != nil {
		p.l.Println("[ERROR] fetching product", err)

		switch err {
		case data.ErrPrdNotFound:
			rw.WriteHeader(http.StatusNotFound)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
			return
		default:
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)
		}

	}
	err = data.ToJSON(product, rw)

	if err != nil {
		p.l.Println("[ERROR] serializing error")
	}

}
