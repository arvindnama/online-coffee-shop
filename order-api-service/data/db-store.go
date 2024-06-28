package data

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arvindnama/golang-microservices/order-service/config"
	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-hclog"
)

type DBOrderStore struct {
	logger hclog.Logger
	DB     *sql.DB
}

func NewDBOrderStore(logger hclog.Logger) (*DBOrderStore, error) {
	cfg := mysql.Config{
		User:                 config.Env.DBUserName,
		Passwd:               config.Env.DBPassword,
		Addr:                 config.Env.DBAddress,
		DBName:               config.Env.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
		MultiStatements:      true,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = initDb(db, logger)
	if err != nil {
		return nil, err
	}

	return &DBOrderStore{logger, db}, nil
}

func initDb(db *sql.DB, logger hclog.Logger) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	logger.Debug("DB: Successfully connected")
	return nil
}

func (store *DBOrderStore) GetAllOrders() ([]*Order, error) {
	sqlStatement := "SELECT * FROM orders JOIN products on orders.id = products.order_id"
	rows, err := store.DB.Query(sqlStatement)

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

	return orders, nil
}

func (store *DBOrderStore) GetOrder(id int64) (*Order, error) {
	sqlStatement := "SELECT * FROM orders JOIN products on orders.id = products.order_id WHERE orders.id=?"
	rows, err := store.DB.Query(sqlStatement, id)

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

	updateStatement := "Update orders SET status = ? where id=?"
	if _, err := tx.ExecContext(ctx, updateStatement, newStatus, id); err != nil {
		return fail(err)
	}

	if err = tx.Commit(); err != nil {
		return fail(err)
	}
	return nil
}

func (store *DBOrderStore) DeleteOrder(ctx context.Context, id int64) error {

	fail := func(err error) error {
		return fmt.Errorf("Delete order status %v", err)
	}

	tx, err := store.DB.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	defer tx.Rollback()

	// First delete all products in the order
	deleteProductsInOrderSQLCmd := "DELETE FROM products where order_id=?"
	if _, err = tx.ExecContext(ctx, deleteProductsInOrderSQLCmd, id); err != nil {
		return fail(err)
	}

	// Second delete the order
	deleteOrderSQLCmd := "DELETE FROM orders where id=?"
	if _, err = tx.ExecContext(ctx, deleteOrderSQLCmd, id); err != nil {
		return fail(err)
	}

	return nil
}

func (store *DBOrderStore) AddOrder(ctx context.Context, order *Order) (int64, error) {

	fail := func(err error) (int64, error) {
		return 0, fmt.Errorf("Create Order: %v", err)
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

	insertIntoOrdersTable := "INSERT INTO orders (name, totalPrice, tx.ExecContext) VALUES(?,?,?)"
	result, err := tx.ExecContext(ctx, insertIntoOrdersTable, order.Name, Initiated, order.TotalPrice)
	if err != nil {
		return fail(err)
	}

	// get the last inserted id from order table
	orderID, err := result.LastInsertId()
	if err != nil {
		return fail(err)
	}

	for _, prod := range order.Products {
		insertIntoProductsTable := "INSERT INTO products (id, order_id, name, quantity, unit_price) VALUES(?,?,?,?, ?)"
		if _, err := tx.ExecContext(ctx, insertIntoProductsTable, prod.ID, orderID, prod.Name, prod.Quantity, prod.UnitPrice); err != nil {
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
	var pName string
	var totalPrice float32
	var status Status

	var product Product = Product{}

	err := rows.Scan(
		&orderId,
		&orderName,
		&totalPrice,
		&status,
		&product.ID,
		&orderId,
		&pName,
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
