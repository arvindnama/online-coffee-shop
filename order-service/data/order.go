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

var ErrOrderNotFound = fmt.Errorf("order not found")

type OrderDB struct {
	orders []*Order
}

var orderDB = &OrderDB{
	orders: []*Order{},
}

func GetAllOrders() []*Order {
	orders := []*Order{}
	for _, o := range orderDB.orders {
		clonedOrder := *o
		orders = append(orders, &clonedOrder)
	}
	return orders
}

func GetOrder(id int64) (*Order, error) {
	idx, err := findOrder(id)

	if idx != -1 {
		return orderDB.orders[idx], nil
	}
	return nil, err
}

func nextOrderId() int64 {
	orderLen := len(orderDB.orders)
	if orderLen == 0 {
		return 1
	}
	lo := orderDB.orders[orderLen-1]
	return lo.ID + 1
}

func UpdateOrder(id int64, order *Order) error {
	order.ID = id
	idx, err := findOrder(order.ID)

	if err != nil {
		return err
	}
	orderDB.orders[idx] = order
	return nil
}

func DeleteOrder(order *Order) error {
	idx, err := findOrder(order.ID)

	if err != nil {
		return err
	}
	orderDB.orders = append(orderDB.orders[:idx], orderDB.orders[idx+1])
	return nil
}

func AddOrder(order *Order) int64 {
	order.ID = nextOrderId()
	orderDB.orders = append(orderDB.orders, order)
	return order.ID
}

func findOrder(id int64) (int, error) {
	for idx, o := range orderDB.orders {
		if o.ID == id {
			return idx, nil
		}
	}
	return -1, ErrOrderNotFound
}

type ValidationError struct {
	Messages []string `json:"messages"`
}
