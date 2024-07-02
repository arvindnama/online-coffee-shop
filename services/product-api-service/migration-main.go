package main

import (
	"os"

	dbUtils "github.com/arvindnama/golang-microservices/libs/utils/db-utils"
	"github.com/arvindnama/golang-microservices/product-api-service/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hashicorp/go-hclog"
)

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "Products Api Service Migration tool",
		Level: hclog.LevelFromString(config.ENV.LogLevel),
	})

	db, err := dbUtils.NewDbConnection(&config.ENV.DBConfig, logger)
	checkDBError("DBConnection", err, logger)

	driver, err := mysql.WithInstance(db, &mysql.Config{})
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
