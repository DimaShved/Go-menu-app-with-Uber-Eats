package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/service"
)

func SetupModifierOptionRoutes(app *fiber.App, modifierOptionService *service.ModifierOptionService, validate *validator.Validate) {
	SetupGenericRoutes[*domain.ModifierOption](app, "/api/modifier-option", modifierOptionService.GenericService(), validate)
}
