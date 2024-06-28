package config

import (
	"fmt"

	envUtils "github.com/arvindnama/golang-microservices/libs/utils/env-utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Address  string
	LogLevel string

	UseDB      bool
	DBAddress  string
	DBUserName string
	DBPassword string
	DBName     string
}

var Env = initConfig()

func initConfig() *Config {
	godotenv.Load()
	return &Config{
		Address:  envUtils.GetEnvString("BIND_ADDRESS", ":9093"),
		LogLevel: envUtils.GetEnvString("LOG_LEVEL", "debug"),

		UseDB: envUtils.GetEnvBool("USE_DB", false),
		DBAddress: fmt.Sprintf(
			"%s:%s",
			envUtils.GetEnvString("DB_HOST", "localhost"),
			envUtils.GetEnvString("DB_PORT", "3306"),
		),
		DBUserName: envUtils.GetEnvString("DB_USERNAME", "root"),
		DBPassword: envUtils.GetEnvString("DB_PASSWORD", ""),
		DBName:     envUtils.GetEnvString("DB_NAME", "baas"),
	}
}
