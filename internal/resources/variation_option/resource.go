package variation_option

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.VariationOption, CreateRequest, UpdateRequest, Response]{
		Name:       "variation_option",
		Path:       "/api/variation-option",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		GetID: func(entity *domain.VariationOption) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.VariationOption, error) {
			return &domain.VariationOption{
				Name:        request.Name,
				Price:       request.Price,
				IsAvailable: request.IsAvailable,
				VariationID: request.VariationID,
			}, nil
		},
		ApplyUpdate: func(entity *domain.VariationOption, request UpdateRequest) error {
			entity.Name = request.Name
			entity.Price = request.Price
			entity.IsAvailable = request.IsAvailable
			entity.VariationID = request.VariationID
			return nil
		},
		MapResponse: mapResponse,
	})
}

func mapResponse(entity *domain.VariationOption) (Response, error) {
	return Response{
		ID:          entity.ID,
		Name:        entity.Name,
		Price:       entity.Price,
		IsAvailable: entity.IsAvailable,
		VariationID: entity.VariationID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}
