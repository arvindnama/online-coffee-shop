package data

import "github.com/hashicorp/go-hclog"

type LocalOrderStore struct {
	orders []*Order
}

func NewLocalOrderStore(logger hclog.Logger) (*LocalOrderStore, error) {
	logger.Debug("localDB: Connected successfully")
	return &LocalOrderStore{
		orders: []*Order{},
	}, nil
}

func (store *LocalOrderStore) GetAllOrders() ([]*Order, error) {
	orders := []*Order{}
	for _, o := range store.orders {
		clonedOrder := *o
		orders = append(orders, &clonedOrder)
	}
	return orders, nil
}

func (store *LocalOrderStore) GetOrder(id int64) (*Order, error) {
	idx, err := store.findOrder(id)

	if idx != -1 {
		return store.orders[idx], nil
	}
	return nil, err
}

func (store *LocalOrderStore) nextOrderId() int64 {
	orderLen := len(store.orders)
	if orderLen == 0 {
		return 1
	}
	lo := store.orders[orderLen-1]
	return lo.ID + 1
}

func (store *LocalOrderStore) UpdateOrderStatus(id int64, newStatus Status) error {
	idx, err := store.findOrder(id)

	if err != nil {
		return err
	}
	store.orders[idx].Status = newStatus
	return nil
}

func (store *LocalOrderStore) DeleteOrder(id int64) error {
	idx, err := store.findOrder(id)

	if err != nil {
		return err
	}
	store.orders = append(store.orders[:idx], store.orders[idx+1])
	return nil
}

func (store *LocalOrderStore) AddOrder(order *Order) (int64, error) {
	order.ID = store.nextOrderId()
	store.orders = append(store.orders, order)
	return order.ID, nil
}

func (store *LocalOrderStore) findOrder(id int64) (int, error) {
	for idx, o := range store.orders {
		if o.ID == id {
			return idx, nil
		}
	}
	return -1, ErrOrderNotFound
}
