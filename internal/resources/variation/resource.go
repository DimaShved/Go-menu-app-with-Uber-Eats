package variation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.Variation, CreateRequest, UpdateRequest, Response]{
		Name:       "variation",
		Path:       "/api/variation",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		Hooks:      Hooks{},
		GetID: func(entity *domain.Variation) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.Variation, error) {
			return &domain.Variation{
				Name:       request.Name,
				CategoryID: request.CategoryID,
			}, nil
		},
		ApplyUpdate: func(entity *domain.Variation, request UpdateRequest) error {
			entity.Name = request.Name
			entity.CategoryID = request.CategoryID
			return nil
		},
		MapResponse: mapResponse,
	})
}

func mapResponse(entity *domain.Variation) (Response, error) {
	options := make([]OptionResponse, 0, len(entity.Options))
	for _, option := range entity.Options {
		options = append(options, mapOptionResponse(option))
	}

	return Response{
		ID:         entity.ID,
		Name:       entity.Name,
		Options:    options,
		CategoryID: entity.CategoryID,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}, nil
}

func mapOptionResponse(entity domain.VariationOption) OptionResponse {
	return OptionResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		Price:       entity.Price,
		IsAvailable: entity.IsAvailable,
		VariationID: entity.VariationID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
