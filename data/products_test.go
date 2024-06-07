package data

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	p := &Product{
		Price: 19,
		SKU:   "abc-abc-def",
	}
	v := NewValidation()
	err := v.Validate(p)
	fmt.Println(err.Errors())
	assert.Len(t, err, 1)

}

func TestProductMissingPriceReturnsErr(t *testing.T) {

	p := &Product{
		Name: "P1",
		SKU:  "abc-cdf-ghi",
	}

	v := NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestProductsToJSON(t *testing.T) {
	ps := []*Product{
		&Product{
			Name: "abc",
		},
	}

	b := bytes.NewBufferString("")
	err := ToJSON(ps, b)
	assert.NoError(t, err)
}
