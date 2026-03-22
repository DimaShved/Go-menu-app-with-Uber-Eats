package menu_item

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func New(db *gorm.DB, validate *validator.Validate) crud.RouteRegistrar {
	return crud.NewHandler(crud.Resource[domain.MenuItem, CreateRequest, UpdateRequest, Response]{
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
		MapResponse: mapResponse,
	})
}

func mapResponse(entity *domain.MenuItem) (Response, error) {
	categories := make([]CategoryResponse, 0, len(entity.Categories))
	for _, category := range entity.Categories {
		categories = append(categories, mapCategoryResponse(category))
	}

	return Response{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Price:       entity.Price,
		IsAvailable: entity.IsAvailable,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Categories:  categories,
	}, nil
}

func mapCategoryResponse(entity domain.MenuCategory) CategoryResponse {
	return CategoryResponse{
		ID:          entity.ID,
		SectionID:   entity.SectionID,
		Name:        entity.Name,
		Description: entity.Description,
		IsAvailable: entity.IsAvailable,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Section:     mapSectionResponse(entity.Section),
	}
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
