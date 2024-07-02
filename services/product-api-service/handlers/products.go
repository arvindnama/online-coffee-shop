package handlers

import (
	"net/http"
	"strconv"

	dataUtils "github.com/arvindnama/golang-microservices/libs/utils/data-utils"
	"github.com/arvindnama/golang-microservices/product-api-service/data"
	"github.com/hashicorp/go-hclog"

	"github.com/gorilla/mux"
)

type Products struct {
	l   hclog.Logger
	v   *dataUtils.Validation
	pDB *data.ProductsStore
}

type KeyProduct struct {
}

func NewProducts(
	l hclog.Logger,
	v *dataUtils.Validation,
	pDB *data.ProductsStore,
) *Products {
	return &Products{l, v, pDB}
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
