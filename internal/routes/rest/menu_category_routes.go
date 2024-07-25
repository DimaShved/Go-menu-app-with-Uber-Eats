package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"uber-go-menu-copy/internal/domain"
	"uber-go-menu-copy/internal/service"
)

func SetupMenuCategoryRoutes(app *fiber.App, menuCategoryService *service.MenuCategoryService, validate *validator.Validate) {
	SetupGenericRoutes[*domain.MenuCategory](app, "/api/menu-category", menuCategoryService, validate)
}
