package data

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"

	"github.com/go-playground/validator/v10"
)

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
	Price float32 `json:"price" validate:"gt=0,required"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
	// json: "-" indicated field will be omitted

}

var ErrPrdNotFound = fmt.Errorf("product not found")

type Products []*Product

func (p *Product) Validate() error {
	validator := validator.New()
	validator.RegisterValidation("sku", validateSKU)
	return validator.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// sku format: xxxx-xxxx-xxxx
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	return len(matches) == 1
}

func (p *Products) ToJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

func AddProduct(prod *Product) {
	prod.ID = getNextId()
	productList = append(productList, prod)
}

func UpdateProduct(id int, prod *Product) error {
	pos, err := findProduct(id)

	if err != nil {
		return err
	}

	prod.ID = id
	productList[pos] = prod
	return nil
}

func DeleteProduct(id int) error {
	pos, err := findProduct(id)

	if err != nil {
		return err
	}

	productList = append(productList[:pos], productList[pos+1])
	return nil
}

func findProduct(id int) (int, error) {
	for i, p := range productList {
		if p.ID == id {
			return i, nil
		}
	}
	return -1, ErrPrdNotFound
}

func getNextId() int {
	lp := productList[len(productList)-1]
	return lp.ID + 1
}

func (p *Product) FromJson(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc-xyz-lmn",
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and string coffee without milk",
		Price:       1.99,
		SKU:         "xyz-ijk-abc",
	},
}

func GetProducts() Products {
	return productList
}
