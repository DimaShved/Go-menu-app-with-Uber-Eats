package menu_item

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.MenuItem, CreateRequest, UpdateRequest, domain.MenuItem]{
		Name:       "menu_item",
		Path:       "/api/menu-item",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		Hooks:      Hooks{},
		GetID: func(entity *domain.MenuItem) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.MenuItem, error) {
			return &domain.MenuItem{
				Name:        request.Name,
				Description: request.Description,
				Price:       request.Price,
				IsAvailable: request.IsAvailable,
			}, nil
		},
		ApplyUpdate: func(entity *domain.MenuItem, request UpdateRequest) error {
			entity.Name = request.Name
			entity.Description = request.Description
			entity.Price = request.Price
			entity.IsAvailable = request.IsAvailable
			return nil
		},
		MapResponse: func(entity *domain.MenuItem) (domain.MenuItem, error) {
			return *entity, nil
		},
	})
}
