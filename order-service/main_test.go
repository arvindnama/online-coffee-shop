package main

import (
	"fmt"
	"testing"

	"github.com/arvindnama/golang-microservices/order-service/sdk/client/orders"
	"github.com/arvindnama/golang-microservices/order-service/sdk/models"
)

func TestOrderClient(t *testing.T) {
	// cfg := client.DefaultTransportConfig().WithHost("localhost:9093")
	// c := client.NewHTTPClientWithConfig(nil, cfg)

	orderName := "order-1"
	order := &models.Order{
		ID:   1,
		Name: &orderName,
		Products: []*models.Product{
			{
				ID:        1,
				Quantity:  10,
				UnitPrice: 1.2,
			},
		},
		TotalPrice: 12,
	}
	client := orders.NewClientWithBearerToken("localhost:9093", "/", "http", "YW5hbWE=")
	params := orders.NewCreateOrderParams().WithBody(order)

	client.CreateOrder(params)
	// _, err := client.CreateOrder(params)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	allOrders, err := client.GetAllOrders(orders.NewGetAllOrdersParams())

	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", allOrders)
}
