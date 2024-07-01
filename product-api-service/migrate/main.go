package main

import (
	"os"

	dbUtils "github.com/arvindnama/golang-microservices/libs/utils/db-utils"
	"github.com/arvindnama/golang-microservices/order-service/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hashicorp/go-hclog"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Products Api Service Migration tool",
		Level: hclog.LevelFromString(config.Env.LogLevel),
	})

	db, err := dbUtils.NewDbConnection(&config.Env.DBConfig, logger)
	checkDBError(err)

	driver, err := mysql.WithInstance(db, &mysql.Config{})
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
