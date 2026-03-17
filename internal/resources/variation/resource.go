package variation

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.Variation, CreateRequest, UpdateRequest, domain.Variation]{
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
		MapResponse: func(entity *domain.Variation) (domain.Variation, error) {
			return *entity, nil
		},
	})
}
