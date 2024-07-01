package data

import "fmt"

// Product defines the structure for an API products
// swagger:model Product
type Product struct {
	//this if of the product
	//
	// required: true
	// min:1
	ID int `json:"id"`

	// the name for this product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float64 `json:"price" validate:"gt=0,required"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
	// json: "-" indicated field will be omitted

}
type Products []*Product

var ErrPrdNotFound = fmt.Errorf("product not found")
