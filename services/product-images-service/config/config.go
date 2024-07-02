package config

import (
	envUtils "github.com/arvindnama/golang-microservices/libs/utils/env-utils"
	"github.com/joho/godotenv"
)

type Config struct {
	Address   string
	LogLevel  string
	ImagePath string
}

var Env = initConfig()

func initConfig() *Config {
	godotenv.Load()
	return &Config{
		Address:   envUtils.GetEnvString("BIND_ADDRESS", ":9091"),
		LogLevel:  envUtils.GetEnvString("LOG_LEVEL", "debug"),
		ImagePath: envUtils.GetEnvString("IMAGE_PATH", "./imagesstore"),
	}
}
