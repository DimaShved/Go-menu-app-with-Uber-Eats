package menu_availability

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

const path = "/api/menu-availability"

type Handler struct {
	crudHandler crud.RouteRegistrar
	repository  *Repository
	txManager   crud.TxManager
	validator   *validator.Validate
}

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	repository := NewRepository()
	txManager := crud.NewTxManager(db)
	crudHandler := crud.NewHandler(crud.Resource[domain.MenuAvailability, CreateRequest, UpdateRequest, domain.MenuAvailability]{
		Name:       "menu_availability",
		Path:       path,
		Repository: repository,
		TxManager:  txManager,
		Validator:  validate,
		GetID: func(entity *domain.MenuAvailability) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.MenuAvailability, error) {
			return &domain.MenuAvailability{
				MenuSectionId: request.MenuSectionID,
				DayOfWeek:     request.DayOfWeek,
				OpenTime:      request.OpenTime,
				CloseTime:     request.CloseTime,
			}, nil
		},
		ApplyUpdate: func(entity *domain.MenuAvailability, request UpdateRequest) error {
			entity.MenuSectionId = request.MenuSectionID
			entity.DayOfWeek = request.DayOfWeek
			entity.OpenTime = request.OpenTime
			entity.CloseTime = request.CloseTime
			return nil
		},
		MapResponse: func(entity *domain.MenuAvailability) (domain.MenuAvailability, error) {
			return *entity, nil
		},
	})

	return &Handler{
		crudHandler: crudHandler,
		repository:  repository,
		txManager:   txManager,
		validator:   validate,
	}
}

func (h *Handler) RegisterRoutes(app *fiber.App) {
	h.crudHandler.RegisterRoutes(app)
	app.Post(path+"/batch", h.batchUpsert)
}

func (h *Handler) batchUpsert(c fiber.Ctx) error {
	var request BatchUpsertRequest
	if err := json.Unmarshal(c.Body(), &request); err != nil {
		return errorx.ErrInvalidInput.WithDetails("Invalid JSON body")
	}

	if h.validator != nil {
		if err := h.validator.Struct(request); err != nil {
			return errorx.ErrInvalidInput.WithDetails(fmt.Sprintf("Input validation failed: %v", err))
		}
	}

	ctx, cancel := context.WithTimeout(c.UserContext(), crud.DefaultTimeout)
	defer cancel()

	response, err := h.upsertBatch(ctx, request)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func (h *Handler) upsertBatch(ctx context.Context, request BatchUpsertRequest) ([]domain.MenuAvailability, error) {
	availabilities, err := buildAvailabilities(request)
	if err != nil {
		return nil, err
	}

	var result []domain.MenuAvailability
	err = h.txManager.WithinTx(ctx, func(tx *gorm.DB) error {
		updated, err := h.repository.UpsertBatch(ctx, tx, availabilities)
		if err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func buildAvailabilities(request BatchUpsertRequest) ([]domain.MenuAvailability, error) {
	seen := make(map[string]struct{}, len(request.Availabilities))
	availabilities := make([]domain.MenuAvailability, 0, len(request.Availabilities))

	for _, availability := range request.Availabilities {
		key := fmt.Sprintf("%s-%d", request.MenuSectionID, availability.DayOfWeek)
		if _, exists := seen[key]; exists {
			return nil, errorx.ErrInvalidInput.WithDetails(
				fmt.Sprintf("duplicate availability for menu section %s on day %d", request.MenuSectionID, availability.DayOfWeek),
			)
		}
		seen[key] = struct{}{}

		availabilities = append(availabilities, domain.MenuAvailability{
			MenuSectionId: request.MenuSectionID,
			DayOfWeek:     availability.DayOfWeek,
			OpenTime:      availability.OpenTime,
			CloseTime:     availability.CloseTime,
		})
	}

	return availabilities, nil
}
