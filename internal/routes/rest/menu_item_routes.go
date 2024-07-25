package rest

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"uber-go-menu-copy/internal/domain"
	"uber-go-menu-copy/internal/inputs"
	"uber-go-menu-copy/internal/service"
)

func SetupMenuItemRoutes(app *fiber.App, menuItemService *service.MenuItemService, validate *validator.Validate) {
	path := "/api/menu-item"

	group := app.Group(path)

	group.Post("/", func(c fiber.Ctx) error {
		var input inputs.MenuItemCreateInput
		if err := json.Unmarshal(c.Body(), &input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
		}

		if err := validate.Struct(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Input validation failed: %v", err.Error())})
		}

		menuItem, err := menuItemService.CreateWithCategories(&domain.MenuItem{
			Name:        input.Name,
			Description: input.Description,
			Price:       input.Price,
			IsAvailable: input.IsAvailable,
		}, input.Categories)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create menu item"})
		}
		return c.Status(fiber.StatusCreated).JSON(menuItem)
	})

	SetupGenericRoutes[*domain.MenuItem](app, path, menuItemService.GenericService(), validate)
}
