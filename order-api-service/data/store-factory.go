package data

import (
	"context"

	"github.com/arvindnama/golang-microservices/order-service/config"
	"github.com/hashicorp/go-hclog"
)

type OrderDatabase interface {
	GetAllOrders() ([]*Order, error)
	AddOrder(ctx context.Context, order *Order) (int64, error)
	GetOrder(id int64) (*Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status Status) error
	DeleteOrder(ctx context.Context, id int64) error
}

func NewOrderStore(logger hclog.Logger) (OrderDatabase, error) {

	if config.Env.UseDB {
		return NewDBOrderStore(logger)
	}
	return NewLocalOrderStore(logger)
}
