package handlers

import (
	"build-go-microservice/data"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
	v *data.Validation
}

type KeyProduct struct {
}

func NewProducts(l *log.Logger, v *data.Validation) *Products {
	return &Products{l, v}
}

//swagger:model
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection error messages
//
//swagger:model
type ValidationError struct {
	Messages []string `json:"messages"`
}

func getProductID(r *http.Request) int {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		panic(err)
	}

	return id
}
