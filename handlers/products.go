package handlers

import (
	"build-go-microservice/data"
	"log"
	"net/http"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		p.getProducts(rw, r)
		return
	case http.MethodPut:
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	products := data.GetProducts()
	err := products.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal produces", http.StatusInternalServerError)
	}

}
