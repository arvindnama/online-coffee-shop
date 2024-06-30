package data

import (
	"context"
	"database/sql"
	"fmt"

	dbUtils "github.com/arvindnama/golang-microservices/libs/utils/db-utils"
	"github.com/arvindnama/golang-microservices/order-service/config"
	"github.com/hashicorp/go-hclog"
)

type DBOrderStore struct {
	logger hclog.Logger
	DB     *sql.DB
}

const (
	GET_ALL_ORDERS_SQL = `
		SELECT orders.id, orders.name,  orders.total_price, orders.status, products.id, products.name, products.quantity, products.unit_price
		FROM (
			SELECT * FROM orders LIMIT ?, ?
		) as orders 
		JOIN products on orders.id = products.order_id
	`
	GET_ORDER_BY_ID_SQL = `
		SELECT orders.id, orders.name,  orders.total_price, orders.status, products.id, products.name, products.quantity, products.unit_price 
		FROM orders 
		JOIN products on orders.id = products.order_id 
		WHERE orders.id=?
	`
	UPDATE_ORDER_STATUS_SQL = `
		UPDATE orders 
		SET status = ? , updated_timestamp=CURRENT_TIMESTAMP
		WHERE id=?
	`
	DELETE_PRODUCTS_IN_ORDER_SQL = `
		DELETE FROM products
		WHERE order_id=?
	`
	DELETE_PRODUCT_SQL = `
		DELETE FROM orders where id=?
	`
	INSERT_ORDER_SQL = `
		INSERT INTO 
		orders (name, status, total_price) 
		VALUES(?,?,?)
	`
	INSERT_PRODUCT_IN_ORDER_SQL = `
		INSERT INTO 
		products (id, order_id, name, quantity, unit_price) 
		VALUES(?,?,?,?, ?)
	`
)

func NewDBOrderStore(logger hclog.Logger) (*DBOrderStore, error) {
	db, err := dbUtils.NewDbConnection(&config.Env.DBConfig, logger)
	return &DBOrderStore{logger, db}, err
}

func (store *DBOrderStore) GetAllOrders(pageNo, pageSize int) ([]*Order, bool, error) {
	// adding one more to the pageSize to evaluate if there is more elements in DB
	// if sql query return `pageSize + 1` entries then there
	// is at-least one more elements present in the db ( in the next page)
	limit := pageSize + 1
	offset := (pageNo - 1) * pageSize

	store.logger.Debug("requested limit", limit, "offset", offset)
	rows, err := store.DB.Query(GET_ALL_ORDERS_SQL, offset, limit)

	if err != nil {
		return nil, false, err
	}

	var orders = []*Order{}
	for rows.Next() {
		orders, err = scanAllOrdersRow(rows, orders)
		if err != nil {
			return nil, false, err
		}
	}
	hasMore := len(orders) == pageSize+1
	return orders, hasMore, nil
}

func (store *DBOrderStore) GetOrder(id int64) (*Order, error) {
	rows, err := store.DB.Query(GET_ORDER_BY_ID_SQL, id)

	if err != nil {
		return nil, err
	}

	var orders = []*Order{}
	for rows.Next() {
		orders, err = scanAllOrdersRow(rows, orders)
		if err != nil {
			return nil, err
		}
	}

	if len(orders) != 1 {
		return nil, ErrOrderNotFound
	}
	return orders[0], nil

}

func (store *DBOrderStore) UpdateOrderStatus(ctx context.Context, id int64, newStatus Status) error {
	fail := func(err error) error {
		return fmt.Errorf("update order status %v", err)
	}

	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	// defer rollback
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, UPDATE_ORDER_STATUS_SQL, newStatus, id); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return nil
}

func (store *DBOrderStore) DeleteOrder(ctx context.Context, id int64) error {

	fail := func(err error) error {
		return fmt.Errorf("delete order status %v", err)
	}

	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	// First delete all products in the order
	if _, err = tx.ExecContext(ctx, DELETE_PRODUCTS_IN_ORDER_SQL, id); err != nil {
		return fail(err)
	}

	// Second delete the order
	if _, err = tx.ExecContext(ctx, DELETE_PRODUCT_SQL, id); err != nil {
		return fail(err)
	}

	return nil
}

func (store *DBOrderStore) AddOrder(ctx context.Context, order *Order) (int64, error) {

	fail := func(err error) (int64, error) {
		return 0, fmt.Errorf("create Order: %v", err)
	}

	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	// Defer a rollback in case anything fails.
	// [learning]:: defer transaction rollback now will call rollback when function exits
	// incase the transaction commit was successful, this rollback operation will result in a no-op
	// in case the transaction fails or commit is not called rollback will be executed.
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, INSERT_ORDER_SQL, order.Name, Initiated, order.TotalPrice)
	if err != nil {
		return fail(err)
	}

	// get the last inserted id from order table
	orderID, err := result.LastInsertId()
	if err != nil {
		return fail(err)
	}

	for _, prod := range order.Products {
		fmt.Printf("Product: %#v\n", prod)
		if _, err := tx.ExecContext(ctx, INSERT_PRODUCT_IN_ORDER_SQL, prod.ID, orderID, prod.Name, prod.Quantity, prod.UnitPrice); err != nil {
			return fail(err)
		}
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fail(err)
	}

	return orderID, nil
}

func scanAllOrdersRow(rows *sql.Rows, orders []*Order) ([]*Order, error) {
	var orderId int
	var orderName string
	var totalPrice float32
	var status Status

	var product Product = Product{}

	err := rows.Scan(
		&orderId,
		&orderName,
		&totalPrice,
		&status,
		&product.ID,
		&product.Name,
		&product.Quantity,
		&product.UnitPrice,
	)

	if err != nil {
		return nil, err
	}
	idx, err := lookupOrder(orders, int64(orderId))

	var order *Order
	if err == ErrOrderNotFound {
		order = &Order{
			ID:         int64(orderId),
			Name:       orderName,
			TotalPrice: float64(totalPrice),
			Status:     status,
			Products:   []*Product{},
		}
		orders = append(orders, order)
	} else {
		order = orders[idx]
	}

	order.Products = append(order.Products, &product)
	return orders, nil
}

func lookupOrder(orders []*Order, id int64) (int, error) {
	for idx, o := range orders {
		if o.ID == id {
			return idx, nil
		}
	}
	return -1, ErrOrderNotFound
}
