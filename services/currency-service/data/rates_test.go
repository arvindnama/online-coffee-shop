package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewExchangeRates(t *testing.T) {

	er, err := NewExchangeRates(hclog.Default())

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", er.rates)
}
