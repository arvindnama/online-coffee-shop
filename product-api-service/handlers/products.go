package handlers

import (
	"log"
	"net/http"
	"strconv"

	currencyClient "github.com/arvindnama/golang-microservices/currency-service/protos"
	"github.com/arvindnama/golang-microservices/product-api-service/data"

	"github.com/gorilla/mux"
)

type Products struct {
	l  *log.Logger
	v  *data.Validation
	cc currencyClient.CurrencyClient
}

type KeyProduct struct {
}

func NewProducts(
	l *log.Logger,
	v *data.Validation,
	cc currencyClient.CurrencyClient,
) *Products {
	return &Products{l, v, cc}
}

//swagger:model GenericError
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection error messages
//
//swagger:model ValidationError
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
