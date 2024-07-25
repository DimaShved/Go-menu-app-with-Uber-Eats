package service

import (
	"fmt"
	"github.com/google/uuid"
	"uber-go-menu-copy/internal/domain"
	"uber-go-menu-copy/internal/repository"
)

type MenuItemService struct {
	genericService IGenericService[*domain.MenuItem]
	menuItemRepo   repository.MenuItemRepo
}

func (s *MenuItemService) GenericService() IGenericService[*domain.MenuItem] {
	return s.genericService
}

func NewMenuItemService(menuItemRepo repository.MenuItemRepo) *MenuItemService {
	genericService := NewGenericService[*domain.MenuItem](menuItemRepo)
	return &MenuItemService{
		genericService: genericService,
		menuItemRepo:   menuItemRepo,
	}
}

func (s *MenuItemService) CreateWithCategories(menuItem *domain.MenuItem, categoryIDs []string) (*domain.MenuItem, error) {
	var uuids []uuid.UUID
	for _, id := range categoryIDs {
		catID, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid category ID: %v", id)
		}
		uuids = append(uuids, catID)
	}

	if err := s.menuItemRepo.CreateWithCategories(menuItem, uuids); err != nil {
		return nil, err
	}

	return menuItem, nil
}
