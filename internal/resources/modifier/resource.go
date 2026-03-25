package modifier

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	repository := NewRepository()
	return crud.NewHandler(crud.Resource[domain.Modifier, CreateRequest, UpdateRequest, Response]{
		Name:       "modifier",
		Path:       "/api/modifier",
		Repository: repository,
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		Hooks:      Hooks{repository: repository},
		MapCreate: func(request CreateRequest) (*domain.Modifier, error) {
			return &domain.Modifier{
				Name:              request.Name,
				TotalMaxSelection: request.TotalMaxSelection,
				CategoryID:        request.CategoryID,
			}, nil
		},
		ApplyUpdate: func(entity *domain.Modifier, request UpdateRequest) error {
			entity.Name = request.Name
			entity.TotalMaxSelection = request.TotalMaxSelection
			entity.CategoryID = request.CategoryID
			return nil
		},
		MapResponse: mapResponse,
	})
}

func mapResponse(entity *domain.Modifier) (Response, error) {
	options := make([]OptionResponse, 0, len(entity.Options))
	for _, option := range entity.Options {
		options = append(options, mapOptionResponse(option))
	}

	return Response{
		ID:                entity.ID,
		Name:              entity.Name,
		TotalMaxSelection: entity.TotalMaxSelection,
		Options:           options,
		CategoryID:        entity.CategoryID,
		CreatedAt:         entity.CreatedAt,
		UpdatedAt:         entity.UpdatedAt,
	}, nil
}

func mapOptionResponse(entity domain.ModifierOption) OptionResponse {
	return OptionResponse{
		ID:           entity.ID,
		Name:         entity.Name,
		Price:        entity.Price,
		MaxSelection: entity.MaxSelection,
		IsAvailable:  entity.IsAvailable,
		ModifierID:   entity.ModifierID,
		CreatedAt:    entity.CreatedAt,
		UpdatedAt:    entity.UpdatedAt,
	}
}
