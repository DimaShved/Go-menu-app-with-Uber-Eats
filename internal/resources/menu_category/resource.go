package menu_category

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.MenuCategory, CreateRequest, UpdateRequest, domain.MenuCategory]{
		Name:       "menu_category",
		Path:       "/api/menu-category",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		GetID: func(entity *domain.MenuCategory) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.MenuCategory, error) {
			return &domain.MenuCategory{
				SectionID:   request.SectionID,
				Name:        request.Name,
				Description: request.Description,
				IsAvailable: request.IsAvailable,
			}, nil
		},
		ApplyUpdate: func(entity *domain.MenuCategory, request UpdateRequest) error {
			entity.SectionID = request.SectionID
			entity.Name = request.Name
			entity.Description = request.Description
			entity.IsAvailable = request.IsAvailable
			return nil
		},
		MapResponse: func(entity *domain.MenuCategory) (domain.MenuCategory, error) {
			return *entity, nil
		},
	})
}
