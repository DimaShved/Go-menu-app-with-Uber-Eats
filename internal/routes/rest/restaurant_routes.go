package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/service"
)

func SetupRestaurantRoutes(app *fiber.App, restaurantService *service.RestaurantService, validate *validator.Validate) {
	SetupGenericRoutes[*domain.Restaurant](app, "/api/restaurants", restaurantService, validate)
}
