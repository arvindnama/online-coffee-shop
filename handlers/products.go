package handlers

import (
	"build-go-microservice/data"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
	case http.MethodPost:
		p.addProduct(rw, r)
		return
	case http.MethodPut:
		idRegex := regexp.MustCompile(`/([0-9]+)`)
		ids := idRegex.FindAllStringSubmatch(r.URL.Path, -1)
		if len(ids) != 1 {
			http.Error(rw, "Invalid URL more than 1 Id", http.StatusBadRequest)
			return
		}
		if len(ids[0]) != 2 {
			http.Error(rw, "Invalid URL more than one capture group", http.StatusBadRequest)
			return
		}
		idString := ids[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(rw, "Invalid ID", http.StatusBadRequest)
		}
		p.updateProduct(id, rw, r)
		return
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)

}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Update product")

	prod, err := readProduct(r.Body)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}

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

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle POST products")

	prod, err := readProduct(r.Body)

	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusBadRequest)
	}

	data.AddProduct(prod)
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle GET products")
	products := data.GetProducts()
	err := products.ToJSON(rw)

	if err != nil {
		http.Error(rw, "Unable to marshal produces", http.StatusInternalServerError)
	}

}

func readProduct(body io.ReadCloser) (*data.Product, error) {
	prod := &data.Product{}
	err := prod.FromJson(body)
	return prod, err
}
