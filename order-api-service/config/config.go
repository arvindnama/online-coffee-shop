package config

import (
	envUtils "github.com/arvindnama/golang-microservices/libs/utils/env-utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Address  string
	LogLevel string
}

var Env = initConfig()

func initConfig() *Config {
	godotenv.Load()
	return &Config{
		Address:  envUtils.GetEnvString("BIND_ADDRESS", ":9093"),
		LogLevel: envUtils.GetEnvString("LOG_LEVEL", "debug"),
	}
}
