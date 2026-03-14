package restaurant

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.Restaurant, CreateRequest, UpdateRequest, domain.Restaurant]{
		Name:       "restaurant",
		Path:       "/api/restaurants",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		GetID: func(entity *domain.Restaurant) uuid.UUID {
			return entity.ID
		},
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
		MapResponse: func(entity *domain.Restaurant) (domain.Restaurant, error) {
			return *entity, nil
		},
	})
}
