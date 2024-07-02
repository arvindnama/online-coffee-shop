package main

import (
	"os"

	"github.com/arvindnama/golang-microservices/order-service/config"
	"github.com/arvindnama/golang-microservices/order-service/data"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hashicorp/go-hclog"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Order Api Service Migration tool",
		Level: hclog.LevelFromString(config.Env.LogLevel),
	})

	store, err := data.NewDBOrderStore(logger)

	checkDBError(err)

	driver, err := mysql.WithInstance(store.DB, &mysql.Config{})
	checkDBError(err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrate/migrations",
		"mysql",
		driver,
	)
	checkDBError(err)

	cmd := os.Args[len(os.Args)-1]

	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			checkDBError(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			checkDBError(err)
		}
	}
}

func checkDBError(err error) {
	if err != nil {
		panic(err)
	}

}
