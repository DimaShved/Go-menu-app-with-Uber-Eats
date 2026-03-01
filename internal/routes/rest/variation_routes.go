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

func SetupVariationRoutes(app *fiber.App, variationService *service.VariationService, validate *validator.Validate) {
	path := "/api/variation"

	group := app.Group(path)

	group.Post("/", func(c fiber.Ctx) error {
		var input inputs.VariationInput
		if err := json.Unmarshal(c.Body(), &input); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(errorx.ErrInvalidInput.Message)
		}

		if err := validate.Struct(input); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(
				errorx.ErrInvalidInput.WithDetails(fmt.Sprintf("Input validation failed: %v", err.Error())),
			)
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		variation, err := variationService.CreateWithOptions(ctx, &input)
		if err != nil {
			return c.Status(errorx.ErrInternal.HTTPStatus).JSON(
				errorx.ErrInternal.WithDetails(fmt.Sprintf("Failed to create menu item: %v", err.Error())),
			)
		}

		return c.Status(fiber.StatusCreated).JSON(variation)
	})

	SetupGenericRoutes[*domain.Variation](app, path, variationService.GenericService(), validate)
}
