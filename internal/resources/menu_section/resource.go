package menu_section

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.MenuSection, CreateRequest, UpdateRequest, Response]{
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
		MapResponse: mapResponse,
	})
}

func mapResponse(entity *domain.MenuSection) (Response, error) {
	response := Response{
		ID:                 entity.ID,
		RestaurantID:       entity.RestaurantID,
		Name:               entity.Name,
		IsAvailable:        entity.IsAvailable,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		Restaurant:         mapRestaurantResponse(entity.Restaurant),
		MenuAvailabilityID: entity.MenuAvailabilityID,
	}
	if entity.MenuAvailability != nil {
		availability := mapMenuAvailabilityResponse(*entity.MenuAvailability)
		response.MenuAvailability = &availability
	}
	return response, nil
}

func mapRestaurantResponse(entity domain.Restaurant) RestaurantResponse {
	return RestaurantResponse{
		ID:        entity.ID,
		Name:      entity.Name,
		Address:   entity.Address,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func mapMenuAvailabilityResponse(entity domain.MenuAvailability) MenuAvailabilityResponse {
	return MenuAvailabilityResponse{
		ID:            entity.ID,
		MenuSectionID: entity.MenuSectionId,
		DayOfWeek:     entity.DayOfWeek,
		OpenTime:      entity.OpenTime,
		CloseTime:     entity.CloseTime,
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
	}
}
