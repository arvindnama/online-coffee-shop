package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "test",
		Price: 19,
		SKU:   "abs-abs-abs",
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
