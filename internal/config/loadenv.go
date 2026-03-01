package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"uber-go-menu/internal/pkg/errorx"
)

func LoadConfig() (Config, error) {
	var config Config

	if err := godotenv.Load(); err != nil {
		return config, errorx.ErrConfigLoad.WithDetails(err.Error())
	}

	if err := envconfig.Process("", &config); err != nil {
		return config, errorx.ErrValidation.WithDetails(err.Error())
	}

	return config, nil
}
