package main

import (
	"github.com/gofiber/fiber/v3"
	"log/slog"
	"os"
	"uber-go-menu-copy/internal/config"
	"uber-go-menu-copy/internal/pkg/db"
	"uber-go-menu-copy/internal/pkg/validator"
	"uber-go-menu-copy/internal/repository"
	"uber-go-menu-copy/internal/routes/rest"
	"uber-go-menu-copy/internal/service"
)

func main() {
	cfg := config.LoadConfig()
	err := db.Connect(&cfg.Database)
	if err != nil {
		os.Exit(1)
	}
	vld := validator.Validate()

	restaurantRepo := repository.NewRestaurantRepo(db.DB)
	restaurantService := service.NewRestaurantService(restaurantRepo)
	menuSectionRepo := repository.NewMenuSectionRepo(db.DB)
	menuSectionService := service.NewMenuSectionService(menuSectionRepo)

	app := fiber.New()

	rest.SetupRestaurantRoutes(app, restaurantService, vld)
	rest.SetupMenuSectionRoutes(app, menuSectionService, vld)

	err = app.Listen(":" + cfg.App.PORT)
	if err != nil {
		slog.Error("Error starting server: %v", err)
	}
}
