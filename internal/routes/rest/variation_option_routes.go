package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/service"
)

func SetupVariationOptionRoutes(app *fiber.App, variationOptionService *service.VariationOptionService, validate *validator.Validate) {
	SetupGenericRoutes[*domain.VariationOption](app, "/api/variation-option", variationOptionService.GenericService(), validate)
}
