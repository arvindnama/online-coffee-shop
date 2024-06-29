package dbUtils

import (
	"database/sql"
	"fmt"

	envUtils "github.com/arvindnama/golang-microservices/libs/utils/env-utils"
	"github.com/go-sql-driver/mysql"
	"github.com/hashicorp/go-hclog"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	UseDB                  bool
	DBAddress              string
	DBUserName             string
	DBPassword             string
	DBName                 string
	DBNet                  string
	DBAllowNativePasswords bool
	DBParseTime            bool
	DBMultiStatements      bool
}

var Env = initDBConfig()

func initDBConfig() *DBConfig {
	godotenv.Load()

	return &DBConfig{

		UseDB: envUtils.GetEnvBool("USE_DB", false),
		DBAddress: fmt.Sprintf(
			"%s:%s",
			envUtils.GetEnvString("DB_HOST", "localhost"),
			envUtils.GetEnvString("DB_PORT", "3306"),
		),
		DBUserName:             envUtils.GetEnvString("DB_USERNAME", "root"),
		DBPassword:             envUtils.GetEnvString("DB_PASSWORD", ""),
		DBName:                 envUtils.GetEnvString("DB_NAME", "baas"),
		DBNet:                  envUtils.GetEnvString("DB_NET", "tcp"),
		DBAllowNativePasswords: envUtils.GetEnvBool("DB_ALLOW_NATIVE_PASSWORDS", true),
		DBParseTime:            envUtils.GetEnvBool("DB_PARSE_TIME", true),
		DBMultiStatements:      envUtils.GetEnvBool("DB_MULTI_STATEMENTS", true),
	}
}

func NewDbConnection(logger hclog.Logger) (*sql.DB, error) {
	config := Env
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
