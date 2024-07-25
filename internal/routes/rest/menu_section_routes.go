package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"uber-go-menu-copy/internal/domain"
	"uber-go-menu-copy/internal/service"
)

func SetupMenuSectionRoutes(app *fiber.App, menuSectionService *service.MenuSectionService, validate *validator.Validate) {
	SetupGenericRoutes[*domain.MenuSections](app, "/api/menu-section", menuSectionService, validate)
}
