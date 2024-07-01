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

import "github.com/arvindnama/golang-microservices/product-api-service/data"

// A list of products
// swagger:response ProductsResponse
type ProductsResponse struct {
	// All products in the system
	// in: body
	Body data.ProductsPaginatedResponse
}

// Product
// swagger:response ProductResponse
type ProductResponse struct {
	// Products
	// in: body
	Body data.Product
}

// swagger:parameters listProduct deleteProduct updateProduct
type ProductIDPathParameterWrapper struct {
	// The id of the product
	// in: path
	// required: true
	ID int `json:"id"`
}

// swagger:response NoContentResponse
type NoContentResponseWrapper struct {
}

// Generic Error message as string
// swagger:response ErrorResponse
type ErrorResponseWrapper struct {
	// Description of Error
	// in: body
	Body GenericError
}

// Validation errors defined as array of string
//
//swagger:response ErrorValidation
type ErrorValidationWrapper struct {
	// collection of validation errors
	//in: body
	Body ValidationError
}

//swagger:parameters listProducts listProduct
type ProductQueryParameter struct {
	// Currency
	//in: query
	//required: false
	Currency string
}

// swagger:parameters listProducts
type ProductsPaginationQueryParameters struct {
	// page no
	// in: query
	// required: false
	// default: 1
	PageNo int `json:"page_no"`

	// page size
	// in:query
	// required: false
	// default: 10
	PageSize int `json:"page_size"`
}
