package data

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

func UpdateOrderStatus(id int64, newStatus Status) error {
	idx, err := findOrder(id)

	if err != nil {
		return err
	}
	orderDB.orders[idx].Status = newStatus
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
