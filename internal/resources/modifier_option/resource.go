package modifier_option

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.ModifierOption, CreateRequest, UpdateRequest, domain.ModifierOption]{
		Name:       "modifier_option",
		Path:       "/api/modifier-option",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		GetID: func(entity *domain.ModifierOption) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.ModifierOption, error) {
			return &domain.ModifierOption{
				Name:         request.Name,
				Price:        request.Price,
				MaxSelection: request.MaxSelection,
				IsAvailable:  request.IsAvailable,
				ModifierID:   request.ModifierID,
			}, nil
		},
		ApplyUpdate: func(entity *domain.ModifierOption, request UpdateRequest) error {
			entity.Name = request.Name
			entity.Price = request.Price
			entity.MaxSelection = request.MaxSelection
			entity.IsAvailable = request.IsAvailable
			entity.ModifierID = request.ModifierID
			return nil
		},
		MapResponse: func(entity *domain.ModifierOption) (domain.ModifierOption, error) {
			return *entity, nil
		},
	})
}
