package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log/slog"
	"os"
)

func LoadConfig() Config {
	var config Config

	if err := godotenv.Load(); err != nil {
		slog.Error(fmt.Sprintf("Error loading .env file: %v", err))
	}

	if err := envconfig.Process("", &config); err != nil {
		slog.Error(fmt.Sprintf("Error config validation: %v", err))
		os.Exit(1)
	}

	return config
}
