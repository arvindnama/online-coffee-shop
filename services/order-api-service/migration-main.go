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

	checkDBError("DBConnection", err, logger)

	driver, err := mysql.WithInstance(store.DB, &mysql.Config{})
	checkDBError("DB Driver creation", err, logger)

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrate/migrations",
		"mysql",
		driver,
	)
	checkDBError("Migration Scripts initialization", err, logger)

	cmd := os.Args[len(os.Args)-1]

	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			checkDBError("Migration up", err, logger)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			checkDBError("Migration down", err, logger)
		}
	}
}

func checkDBError(source string, err error, logger hclog.Logger) {
	if err != nil {
		logger.Error(source, err)
		panic(err)
	}

}
