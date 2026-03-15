package menu_availability

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	repository := NewRepository()
	txManager := crud.NewTxManager(db)

	return crud.NewHandler(crud.Resource[domain.MenuAvailability, CreateRequest, UpdateRequest, domain.MenuAvailability]{
		Name:       "menu_availability",
		Path:       "/api/menu-availability",
		Repository: repository,
		TxManager:  txManager,
		Validator:  validate,
		GetID: func(entity *domain.MenuAvailability) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.MenuAvailability, error) {
			return nil, errorx.ErrInvalidInput.WithDetails("menu availability create is handled as a batch upsert")
		},
		CreateOverride: func(ctx context.Context, request CreateRequest) (any, error) {
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
			return result, nil
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
}

func buildAvailabilities(request CreateRequest) ([]domain.MenuAvailability, error) {
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
