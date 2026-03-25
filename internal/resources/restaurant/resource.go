package restaurant

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.Restaurant, CreateRequest, UpdateRequest, Response]{
		Name:       "restaurant",
		Path:       "/api/restaurants",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		MapCreate: func(request CreateRequest) (*domain.Restaurant, error) {
			return &domain.Restaurant{
				Name:    request.Name,
				Address: request.Address,
			}, nil
		},
		ApplyUpdate: func(entity *domain.Restaurant, request UpdateRequest) error {
			entity.Name = request.Name
			entity.Address = request.Address
			return nil
		},
		MapResponse: mapResponse,
	})
}

func mapResponse(entity *domain.Restaurant) (Response, error) {
	return Response{
		ID:        entity.ID,
		Name:      entity.Name,
		Address:   entity.Address,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}
