package crud

import (
	"context"
	"encoding/json"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"uber-go-menu/internal/pkg/errorx"
)

type ResourceHandler[Entity any, CreateRequest any, UpdateRequest any, Response any] struct {
	resource Resource[Entity, CreateRequest, UpdateRequest, Response]
	service  *ResourceService[Entity, CreateRequest, UpdateRequest, Response]
}

func NewHandler[Entity any, CreateRequest any, UpdateRequest any, Response any](resource Resource[Entity, CreateRequest, UpdateRequest, Response]) *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response] {
	resource = resource.prepareForHandler()
	return &ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]{
		resource: resource,
		service:  NewService(resource),
	}
}

func (h *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]) RegisterRoutes(app *fiber.App) {
	group := app.Group(h.resource.Path)

	group.Post("/", h.create)
	group.Get("/", h.list)
	group.Get("/:id", h.getByID)
	group.Put("/:id", h.update)
	group.Delete("/:id", h.delete)
}

func (h *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]) create(c fiber.Ctx) error {
	request, err := decodeBody[CreateRequest](c)
	if err != nil {
		return err
	}

	ctx, cancel := h.requestContext(c)
	defer cancel()

	response, err := h.service.Create(ctx, request)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(response)
}

func (h *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]) list(c fiber.Ctx) error {
	ctx, cancel := h.requestContext(c)
	defer cancel()

	response, err := h.service.List(ctx)
	if err != nil {
		return err
	}
	return c.JSON(response)
}

func (h *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]) getByID(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	ctx, cancel := h.requestContext(c)
	defer cancel()

	response, err := h.service.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return c.JSON(response)
}

func (h *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]) update(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	request, err := decodeBody[UpdateRequest](c)
	if err != nil {
		return err
	}

	ctx, cancel := h.requestContext(c)
	defer cancel()

	response, err := h.service.Update(ctx, id, request)
	if err != nil {
		return err
	}
	return c.JSON(response)
}

func (h *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]) delete(c fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return err
	}

	ctx, cancel := h.requestContext(c)
	defer cancel()

	if err := h.service.Delete(ctx, id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *ResourceHandler[Entity, CreateRequest, UpdateRequest, Response]) requestContext(c fiber.Ctx) (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.UserContext(), h.resource.Timeout)
}

func decodeBody[Request any](c fiber.Ctx) (Request, error) {
	var request Request
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return request, errorx.ErrInvalidInput.WithDetails("Invalid JSON body")
	}
	return request, nil
}

func parseID(c fiber.Ctx) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return uuid.Nil, errorx.ErrInvalidInput.WithDetails("Invalid UUID format")
	}
	return id, nil
}
