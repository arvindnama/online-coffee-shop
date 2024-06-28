// Package classification for Order API
//
// # Documentation for Order API
//
// Schemes: Http
// BasePath: /
// Version: 1.0.0
//
// Consumes:
//   - application/json
//
// Produces:
//   - application/json
//
// swagger:meta
package handler

import "github.com/arvindnama/golang-microservices/order-service/data"

// A orders
// swagger:parameters createOrder patchOrder
type OrderRequest struct {
	// A order in DB
	// in: body
	Body data.Order
}

// A orders
// swagger:response OrderResponse
type OrderResponse struct {
	// A order in DB
	// in: body
	Body data.Order
}

// A list of orders
// swagger:response OrdersResponse
type OrdersResponse struct {
	// All orders in DB
	// in: body
	Body []data.Order
}

// swagger:parameters getOrder patchOrder
type OrderIdPathParameter struct {
	// Id of the product
	// in: path
	// required: true
	ID int64 `json:"id"`
}

// swagger:response ErrorResponse
type ErrorResponse struct {
	// collection of errors
	//in: body
	Body data.ValidationError
}
