// Package classification of Product API
//
// Documentation for Product API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import "build-go-microservice/data"

// A list of products
// swagger:response productsResponse
type ProductsResponse struct {
	// All products in the system
	// in: body
	Body []data.Products
}

// swagger:parameters deleteProduct
type productIDPathParameterWrapper struct {
	// The id of the product to delete
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// Generic Error message as string
// swagger:response errorResponse
type ErrorResponseWrapper struct {
	// Description of Error
	// in: body
	Body GenericError
}

// Validation errors defined as array of string
//
//swagger:response errorValidation
type ErrorValidationWrapper struct {
	// collection of validation errors
	//in: body
	Body ValidationError
}
