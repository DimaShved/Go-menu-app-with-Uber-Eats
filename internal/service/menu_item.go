package service

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/repository"
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

func (s *MenuItemService) CreateWithCategories(ctx context.Context, menuItem *domain.MenuItem, categoryIDs []string) (*domain.MenuItem, error) {
	uuids := make([]uuid.UUID, 0, len(categoryIDs))
	for _, id := range categoryIDs {
		catID, err := uuid.Parse(id)
		if err != nil {
			return nil, errorx.NewAPIError(107, http.StatusBadRequest, "invalid category ID: "+id, "")
		}
		uuids = append(uuids, catID)
	}

	err := s.menuItemRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.menuItemRepo.CreateWithCategoriesTx(ctx, tx, menuItem, uuids)
	})
	if err != nil {
		return nil, err
	}

	return menuItem, nil
}
