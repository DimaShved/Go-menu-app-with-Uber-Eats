package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"time"
	"uber-go-menu/internal/domain/interfaces"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/service"
)

func SetupGenericRoutes[T interfaces.Identifiable](app *fiber.App, path string, service service.IGenericService[T], validate *validator.Validate) {
	group := app.Group(path)

	group.Post("/", func(c fiber.Ctx) error {
		var entity T
		if err := json.Unmarshal(c.Body(), &entity); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(errorx.ErrInvalidInput.Message)
		}

		if identifiable, ok := any(entity).(interfaces.Identifiable); ok {
			if identifiable.GetID() != uuid.Nil {
				return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(
					errorx.ErrInvalidInput.WithDetails("Creating new entity with ID is not allowed"),
				)
			}
		}

		if err := validate.Struct(entity); err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(
				errorx.ErrInvalidInput.WithDetails(fmt.Sprintf("Input validation failed: %v", err.Error())),
			)
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		result, err := service.Save(ctx, entity)
		if err != nil {
			return c.Status(errorx.ErrInternal.HTTPStatus).JSON(errorx.ErrInternal.WithDetails(err.Error()))
		}
		return c.Status(fiber.StatusCreated).JSON(result)
	})

	group.Get("/", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		entities, err := service.FindAll(ctx)
		if err != nil {
			return c.Status(errorx.ErrInternal.HTTPStatus).JSON(errorx.ErrInternal.WithDetails(err.Error()))
		}
		return c.JSON(entities)
	})

	group.Get("/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(
				errorx.ErrInvalidInput.WithDetails("Invalid UUID format"),
			)
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		entity, err := service.FindOneById(ctx, id)
		if err != nil {
			return c.Status(errorx.ErrNotFound.HTTPStatus).JSON(
				errorx.ErrNotFound.WithDetails("Entity not found"),
			)
		}
		return c.JSON(entity)
	})

	group.Delete("/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(errorx.ErrInvalidInput.HTTPStatus).JSON(
				errorx.ErrInvalidInput.WithDetails("Invalid UUID format"),
			)
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
		defer cancel()

		err = service.SoftDelete(ctx, id)
		if err != nil {
			return c.Status(errorx.ErrInternal.HTTPStatus).JSON(
				errorx.ErrInternal.WithDetails(fmt.Sprintf("Failed to delete entity: %v", err.Error())),
			)
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
}
