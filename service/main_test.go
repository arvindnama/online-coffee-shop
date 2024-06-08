package main

import (
	"build-go-microservice/sdk/client"
	"build-go-microservice/sdk/client/products"
	"fmt"
	"testing"
)

func TestProductClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:9090")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()

	prod, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("%#v", prod.GetPayload()[0])
}
