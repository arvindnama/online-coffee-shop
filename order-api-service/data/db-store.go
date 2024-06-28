package data

import (
	"database/sql"

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
	return nil, nil
}

func (store *DBOrderStore) GetOrder(id int64) (*Order, error) {
	return nil, nil

}

func (store *DBOrderStore) UpdateOrderStatus(id int64, newStatus Status) error {
	return nil
}

func (store *DBOrderStore) DeleteOrder(id int64) error {
	return nil
}

func (store *DBOrderStore) AddOrder(order *Order) (int64, error) {
	return -1, nil
}
