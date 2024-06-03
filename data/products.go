package data

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type Product struct {
	ID          int     `json:"id"`   // These are called field tags (annotations)
	Name        string  `json:"name"` // json:"<>" we can change the name of the field (while marshalling)
	Description string  `json:"description"`
	Price       float32 `json:"price"`
	SKU         string  `json:"sku"`
	CreatedOn   string  `json:"-"` // json: "-" indicated field will be omitted
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

var ErrPrdNotFound = fmt.Errorf("product not found")

type Products []*Product

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
		SKU:         "abc323",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and string coffee without milk",
		Price:       1.99,
		SKU:         "fcf123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}

func GetProducts() Products {
	return productList
}
