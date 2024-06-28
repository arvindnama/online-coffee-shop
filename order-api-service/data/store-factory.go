package data

type OrderDatabase interface {
	GetAllOrders() ([]*Order, error)
	AddOrder(order *Order) (int64, error)
	GetOrder(id int64) (*Order, error)
	UpdateOrderStatus(id int64, status Status) error
	DeleteOrder(id int64) error
}

func NewOrderStore() OrderDatabase {
	return NewLocalOrderStore()
}
