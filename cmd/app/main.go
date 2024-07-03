package main

import (
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"os"
	"uber-go-menu-copy/internal/config"
	"uber-go-menu-copy/internal/pkg/db"
	"uber-go-menu-copy/internal/pkg/validator"
)

func main() {
	cfg := config.LoadConfig()
	err := db.Connect(&cfg.Database)
	if err != nil {
		os.Exit(1)
	}
	_ = validator.Validate()

	app := fiber.New()

	err = app.Listen(":" + cfg.App.PORT)
	if err != nil {
		slog.Error("Error starting server: %v", err)
	}
}
