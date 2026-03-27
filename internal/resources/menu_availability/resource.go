package menu_availability

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

const path = "/api/menu-availability"

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	repository := NewRepository()
	txManager := crud.NewTxManager(db)
	return crud.NewHandler(crud.Resource[domain.MenuAvailability, CreateRequest, UpdateRequest, Response]{
		Name:       "menu_availability",
		Path:       path,
		Repository: repository,
		TxManager:  txManager,
		Validator:  validate,
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
		MapResponse: mapResponse,
		ExtraRoutes: func(router fiber.Router) {
			router.Post("/batch", batchUpsert(repository, txManager, validate))
		},
	})
}

func batchUpsert(repository *Repository, txManager crud.TxManager, validate *validator.Validate) fiber.Handler {
	return func(c fiber.Ctx) error {
		var request BatchUpsertRequest
		if err := json.Unmarshal(c.Body(), &request); err != nil {
			return errorx.ErrInvalidInput.WithDetails("Invalid JSON body")
		}

		if validate != nil {
			if err := validate.Struct(request); err != nil {
				return errorx.ErrInvalidInput.WithDetails(fmt.Sprintf("Input validation failed: %v", err))
			}
		}

		ctx, cancel := context.WithTimeout(c.UserContext(), crud.DefaultTimeout)
		defer cancel()

		response, err := upsertBatch(ctx, repository, txManager, request)
		if err != nil {
			return err
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}

func upsertBatch(ctx context.Context, repository *Repository, txManager crud.TxManager, request BatchUpsertRequest) ([]Response, error) {
	availabilities, err := buildAvailabilities(request)
	if err != nil {
		return nil, err
	}

	var result []domain.MenuAvailability
	err = txManager.WithinTx(ctx, func(tx *gorm.DB) error {
		updated, err := repository.UpsertBatch(ctx, tx, availabilities)
		if err != nil {
			return err
		}
		result = updated
		return nil
	})
	if err != nil {
		return nil, err
	}

	response := make([]Response, 0, len(result))
	for i := range result {
		item, err := mapResponse(&result[i])
		if err != nil {
			return nil, err
		}
		response = append(response, item)
	}
	return response, nil
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

func mapResponse(entity *domain.MenuAvailability) (Response, error) {
	return Response{
		ID:            entity.ID,
		MenuSectionID: entity.MenuSectionId,
		DayOfWeek:     entity.DayOfWeek,
		OpenTime:      entity.OpenTime,
		CloseTime:     entity.CloseTime,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}, nil
}
