package data

import (
	"fmt"
)

type Product struct {
	ID        int64   `json:"id" validate:"required"`
	Quantity  int64   `json:"quantity" validate:"required"`
	UnitPrice float32 `json:"unitPrice" validate:"required"`
}

type Status string

const (
	Initiated  Status = "initiated"
	Processing Status = "processing"
	Cancelled  Status = "cancelled"
	Completed  Status = "completed"
)

type Order struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name" validate:"required"`
	Products   []Product `json:"products" validate:"required"`
	TotalPrice float64   `json:"totalPrice"`
	Status     Status    `json:"status"`
}

type ValidationError struct {
	Messages []string `json:"messages"`
}

var ErrOrderNotFound = fmt.Errorf("order not found")
