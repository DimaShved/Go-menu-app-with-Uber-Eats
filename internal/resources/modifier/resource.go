package modifier

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.Modifier, CreateRequest, UpdateRequest, domain.Modifier]{
		Name:       "modifier",
		Path:       "/api/modifier",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		Hooks:      Hooks{},
		GetID: func(entity *domain.Modifier) uuid.UUID {
			return entity.ID
		},
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
		MapResponse: func(entity *domain.Modifier) (domain.Modifier, error) {
			return *entity, nil
		},
	})
}
