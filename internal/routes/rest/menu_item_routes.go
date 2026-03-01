package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"time"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/inputs"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/service"
)

func SetupMenuItemRoutes(app *fiber.App, menuItemService *service.MenuItemService, validate *validator.Validate) {
	path := "/api/menu-item"

	group := app.Group(path)

	group.Post("/", func(c fiber.Ctx) error {
		var input inputs.MenuItemCreateInput
		if err := json.Unmarshal(c.Body(), &input); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(errorx.ErrInvalidInput.Message)
		}

		if err := validate.Struct(input); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(errorx.ErrInvalidInput.WithDetails(fmt.Sprintf("Input validation failed: %v", err.Error())))
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		menuItem, err := menuItemService.CreateWithCategories(ctx, &domain.MenuItem{
			Name:        input.Name,
			Description: input.Description,
			Price:       input.Price,
			IsAvailable: input.IsAvailable,
		}, input.Categories)
		if err != nil {
			return c.Status(errorx.ErrInternal.HTTPStatus).JSON(errorx.ErrInternal.WithDetails(fmt.Sprintf("Failed to create menu item: %v", err.Error())))
		}
		return c.Status(fiber.StatusCreated).JSON(menuItem)
	})

	SetupGenericRoutes[*domain.MenuItem](app, path, menuItemService.GenericService(), validate)
}
