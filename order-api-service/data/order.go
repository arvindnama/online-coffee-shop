package data

import (
	"fmt"
)

//swagger:model Product
type Product struct {
	// Identifier of the product in the order
	//	require:true
	ID int64 `json:"id" validate:"required"`

	// quantity of products purchased
	//	require:true
	Quantity int64 `json:"quantity" validate:"required"`

	// price of one product
	//	require:true
	UnitPrice float32 `json:"unitPrice" validate:"required"`
}

type Status string

const (
	Initiated  Status = "initiated"
	Processing Status = "processing"
	Cancelled  Status = "cancelled"
	Completed  Status = "completed"
)

//swagger:model Order
type Order struct {
	// Identifier of the order
	ID int64 `json:"id"`
	// name of the order
	// required: true
	Name string `json:"name" validate:"required"`
	// products purchased in the order
	// required: true
	Products []*Product `json:"products" validate:"required"`
	// total cost of the order
	TotalPrice float64 `json:"totalPrice"`
	// order status
	Status Status `json:"status"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

var ErrOrderNotFound = fmt.Errorf("order not found")
