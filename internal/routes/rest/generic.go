package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"time"
	"uber-go-menu-copy/internal/domain/interfaces"
	"uber-go-menu-copy/internal/service"
)

func SetupGenericRoutes[T any](app *fiber.App, path string, service service.IGenericService[T], validate *validator.Validate) {
	group := app.Group(path)

	group.Post("/", func(c fiber.Ctx) error {
		var entity T
		if err := json.Unmarshal(c.Body(), &entity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request data"})
		}

		if identifiable, ok := any(entity).(interfaces.Identifiable); ok {
			if identifiable.GetID() != uuid.Nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Creating new entity with ID is not allowed"})
			}
		}

		if err := validate.Struct(entity); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Input validation failed: %v", err.Error())})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		result, err := service.Save(ctx, entity)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create entity"})
		}
		return c.Status(fiber.StatusCreated).JSON(result)
	})

	group.Get("/", func(c fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		entities, err := service.FindAll(ctx)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve entities"})
		}
		return c.JSON(entities)
	})

	group.Get("/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid UUID format"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		entity, err := service.FindOneById(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Entity not found"})
		}
		return c.JSON(entity)
	})

	group.Delete("/:id", func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := uuid.Parse(idStr)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = service.SoftDelete(ctx, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete entity"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	})
}
