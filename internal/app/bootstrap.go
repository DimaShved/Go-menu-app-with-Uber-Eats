package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"uber-go-menu/internal/platform/crud"
	"uber-go-menu/internal/resources/menu_availability"
	"uber-go-menu/internal/resources/menu_category"
	"uber-go-menu/internal/resources/menu_item"
	"uber-go-menu/internal/resources/menu_section"
	"uber-go-menu/internal/resources/modifier"
	"uber-go-menu/internal/resources/modifier_option"
	"uber-go-menu/internal/resources/restaurant"
	"uber-go-menu/internal/resources/variation"
	"uber-go-menu/internal/resources/variation_option"
)

func NewHTTPServer(database *gorm.DB, vld *validator.Validate) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})

	app.Use(requestLogger())

	registry := crud.NewRegistry()
	registry.Register(restaurant.New(database, vld))
	registry.Register(menu_section.New(database, vld))
	registry.Register(menu_category.New(database, vld))
	registry.Register(menu_item.New(database, vld))
	registry.Register(menu_availability.New(database, vld))
	registry.Register(variation.New(database, vld))
	registry.Register(variation_option.New(database, vld))
	registry.Register(modifier.New(database, vld))
	registry.Register(modifier_option.New(database, vld))
	registry.Mount(app)

	return app
}
