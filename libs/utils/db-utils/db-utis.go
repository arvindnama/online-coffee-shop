package dbUtils

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-hclog"
)

func NewDbConnection(config *DBConfig, logger hclog.Logger) (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 config.DBUserName,
		Passwd:               config.DBPassword,
		Addr:                 config.DBAddress,
		DBName:               config.DBName,
		Net:                  config.DBNet,
		AllowNativePasswords: config.DBAllowNativePasswords,
		ParseTime:            config.DBParseTime,
		MultiStatements:      config.DBMultiStatements,
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	err = initDb(db, logger)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initDb(db *sql.DB, logger hclog.Logger) error {
	err := db.Ping()
	if err != nil {
		return err
	}
	logger.Debug("DB: Successfully connected")
	return nil
}
