package menu_category

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.MenuCategory, CreateRequest, UpdateRequest, Response]{
		Name:       "menu_category",
		Path:       "/api/menu-category",
		Repository: NewRepository(),
		TxManager:  crud.NewTxManager(db),
		Validator:  validate,
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
		MapResponse: mapResponse,
	})
}

func mapResponse(entity *domain.MenuCategory) (Response, error) {
	items := make([]MenuItemSummary, 0, len(entity.Items))
	for _, item := range entity.Items {
		items = append(items, mapMenuItemSummary(item))
	}

	return Response{
		ID:          entity.ID,
		SectionID:   entity.SectionID,
		Name:        entity.Name,
		Description: entity.Description,
		IsAvailable: entity.IsAvailable,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Section:     mapSectionResponse(entity.Section),
		Items:       items,
	}, nil
}

func mapSectionResponse(entity domain.MenuSection) SectionResponse {
	return SectionResponse{
		ID:                 entity.ID,
		RestaurantID:       entity.RestaurantID,
		Name:               entity.Name,
		IsAvailable:        entity.IsAvailable,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		Restaurant:         mapRestaurantResponse(entity.Restaurant),
		MenuAvailabilityID: entity.MenuAvailabilityID,
	}
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

func mapMenuItemSummary(entity domain.MenuItem) MenuItemSummary {
	return MenuItemSummary{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
		IsAvailable: entity.IsAvailable,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}
}
