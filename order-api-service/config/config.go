package config

import (
	dbUtils "github.com/arvindnama/golang-microservices/libs/utils/db-utils"
	envUtils "github.com/arvindnama/golang-microservices/libs/utils/env-utils"
	"github.com/joho/godotenv"
)

type Config struct {
	dbUtils.DBConfig
	Address  string
	LogLevel string
}

var Env = initConfig()

func initConfig() *Config {
	godotenv.Load()

	return &Config{
		DBConfig: *dbUtils.Env,

		Address:  envUtils.GetEnvString("BIND_ADDRESS", ":9093"),
		LogLevel: envUtils.GetEnvString("LOG_LEVEL", "debug"),
	}
}
