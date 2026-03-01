package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/domain/interfaces"
	"uber-go-menu/internal/pkg/errorx"
)

type MenuItemRepo interface {
	IGenericRepository[*domain.MenuItem]
	DB() *gorm.DB
	CreateWithCategories(menuItem *domain.MenuItem, categoryIDs []uuid.UUID) error
	CreateWithCategoriesTx(ctx context.Context, tx *gorm.DB, menuItem *domain.MenuItem, categoryIDs []uuid.UUID) error
}

type menuItemRepo struct {
	*GenericRepository[*domain.MenuItem]
}

func NewMenuItemRepo(db *gorm.DB) MenuItemRepo {
	return &menuItemRepo{NewGenericRepo[*domain.MenuItem](db)}
}

func (r *menuItemRepo) DB() *gorm.DB {
	return r.db
}

func (r *menuItemRepo) CreateWithCategories(menuItem *domain.MenuItem, categoryIDs []uuid.UUID) error {
	return r.CreateWithCategoriesTx(context.Background(), r.db, menuItem, categoryIDs)
}

func (r *menuItemRepo) CreateWithCategoriesTx(ctx context.Context, tx *gorm.DB, menuItem *domain.MenuItem, categoryIDs []uuid.UUID) error {
	if err := tx.WithContext(ctx).Create(menuItem).Error; err != nil {
		return _handleDbError(err, "CreateMenuItem")
	}

	for _, catID := range categoryIDs {
		var category domain.MenuCategory
		if err := tx.WithContext(ctx).First(&category, "id = ?", catID).Error; err != nil {
			return errorx.ErrDatabaseQuery.WithDetails(
				fmt.Sprintf("failed to find category with ID %v: %v", catID, err),
			)
		}

		if err := tx.WithContext(ctx).Model(menuItem).Association("Categories").Append(&category); err != nil {
			return errorx.ErrDatabaseQuery.WithDetails(
				fmt.Sprintf("failed to associate category %v with menu item: %v", catID, err),
			)
		}
	}

	if identifiable, ok := any(menuItem).(interfaces.Identifiable); ok {
		query := tx.WithContext(ctx).Model(&domain.MenuItem{}).Where("id = ?", identifiable.GetID())
		if preloader, ok := any(menuItem).(interfaces.Preloader); ok {
			for _, relation := range preloader.PreloadRelations() {
				query = query.Preload(relation)
			}
		}
		if err := query.First(menuItem).Error; err != nil {
			return _handleDbError(err, "ReloadMenuItem")
		}
		return nil
	}

	return errorx.NewAppError(300, "MenuItem does not implement Identifiable interface")
}
