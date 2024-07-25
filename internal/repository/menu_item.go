package repository

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
	"uber-go-menu-copy/internal/domain"
	"uber-go-menu-copy/internal/domain/interfaces"
)

type MenuItemRepo interface {
	IGenericRepository[*domain.MenuItem]
	CreateWithCategories(menuItem *domain.MenuItem, categoryIDs []uuid.UUID) error
}

type menuItemRepo struct {
	*GenericRepository[*domain.MenuItem]
}

func NewMenuItemRepo(db *gorm.DB) MenuItemRepo {
	return &menuItemRepo{NewGenericRepo[*domain.MenuItem](db)}
}

func (r *menuItemRepo) CreateWithCategories(menuItem *domain.MenuItem, categoryIDs []uuid.UUID) error {
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Create(menuItem).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, catID := range categoryIDs {
		var category domain.MenuCategory
		if err := tx.First(&category, "id = ?", catID).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to find category with ID %v: %w", catID, err)
		}
		if err := tx.Model(menuItem).Association("Categories").Append(&category); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to associate category %v with menu item: %w", catID, err)
		}
	}

	if identifiable, ok := any(menuItem).(interfaces.Identifiable); ok {
		query := tx.Model(&domain.MenuItem{}).Where("id = ?", identifiable.GetID())
		if preloader, ok := any(menuItem).(interfaces.Preloader); ok {
			for _, relation := range preloader.PreloadRelations() {
				query = query.Preload(relation)
			}
		}
		if err := query.First(menuItem).Error; err != nil {
			tx.Rollback()
			return err
		}
	} else {
		slog.Error("MenuItem does not implement Identifiable interface")
		tx.Rollback()
		return fmt.Errorf("MenuItem does not implement Identifiable interface")
	}

	return tx.Commit().Error
}
