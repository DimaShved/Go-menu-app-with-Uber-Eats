package rest

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/inputs"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/service"
)

func SetupMenuAvailabilityRoutes(app *fiber.App, menuAvailabilityService *service.MenuAvailabilityService, validate *validator.Validate) {
	path := "/api/menu-availability"

	group := app.Group(path)

	group.Post("/", func(c fiber.Ctx) error {
		var input inputs.MenuAvailabilityInput
		if err := json.Unmarshal(c.Body(), &input); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(errorx.ErrInvalidInput.Message)
		}

		if err := validate.Struct(input); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(errorx.ErrInvalidInput.WithDetails(fmt.Sprintf("Input validation failed: %v", err.Error())))
		}

		availabilities := make([]*domain.MenuAvailability, len(input.Availabilities))
		for i, a := range input.Availabilities {
			availabilities[i] = &domain.MenuAvailability{
				MenuSectionId: input.MenuSectionId,
				DayOfWeek:     a.DayOfWeek,
				OpenTime:      a.OpenTime,
				CloseTime:     a.CloseTime,
			}
		}

		updatedAvailabilities, err := menuAvailabilityService.UpsertBatch(availabilities)
		if err != nil {
			return c.Status(errorx.ErrInternal.HTTPStatus).JSON(errorx.ErrInternal.WithDetails(fmt.Sprintf("Failed to create or update menu availability: %v", err.Error())))
		}

		return c.Status(fiber.StatusCreated).JSON(updatedAvailabilities)
	})

	SetupGenericRoutes[*domain.MenuAvailability](app, path, menuAvailabilityService.GenericService(), validate)
}
