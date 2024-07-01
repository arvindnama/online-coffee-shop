package dbUtils

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-hclog"
	gormMysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDbConnection(config *DBConfig, logger hclog.Logger) (*sql.DB, error) {

	db, err := createSqlConnection(config)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}
	logger.Debug("DB: Successfully connected")

	if err != nil {
		return nil, err
	}

	return db, nil
}

func createSqlConnection(config *DBConfig) (*sql.DB, error) {
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
	return sql.Open("mysql", cfg.FormatDSN())
}

func NewGormDbConnection(config *DBConfig, logger hclog.Logger) (*gorm.DB, error) {

	sqlDB, err := createSqlConnection(config)

	if err != nil {
		return nil, err
	}
	return gorm.Open(gormMysqlDriver.New(gormMysqlDriver.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
}
