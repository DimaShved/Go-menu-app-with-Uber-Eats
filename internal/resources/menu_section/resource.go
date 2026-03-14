package menu_section

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.MenuSection, CreateRequest, UpdateRequest, domain.MenuSection]{
		Name:       "menu_section",
		Path:       "/api/menu-section",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
		GetID: func(entity *domain.MenuSection) uuid.UUID {
			return entity.ID
		},
		MapCreate: func(request CreateRequest) (*domain.MenuSection, error) {
			return &domain.MenuSection{
				RestaurantID:       request.RestaurantID,
				Name:               request.Name,
				IsAvailable:        request.IsAvailable,
				MenuAvailabilityID: request.MenuAvailabilityID,
			}, nil
		},
		ApplyUpdate: func(entity *domain.MenuSection, request UpdateRequest) error {
			entity.RestaurantID = request.RestaurantID
			entity.Name = request.Name
			entity.IsAvailable = request.IsAvailable
			entity.MenuAvailabilityID = request.MenuAvailabilityID
			return nil
		},
		MapResponse: func(entity *domain.MenuSection) (domain.MenuSection, error) {
			return *entity, nil
		},
	})
}
