package menu_item

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.MenuItem]
}

func NewRepository() *Repository {
	return &Repository{
		GormRepository: crud.NewGormRepository[domain.MenuItem](crud.QueryOptions{
			Preloads: []string{"Categories", "Categories.Section", "Categories.Section.Restaurant"},
		}),
	}
}

func (r *Repository) AttachCategories(ctx context.Context, tx *gorm.DB, entity *domain.MenuItem, categoryIDs []uuid.UUID) error {
	for _, categoryID := range categoryIDs {
		var category domain.MenuCategory
		if err := tx.WithContext(ctx).First(&category, "id = ?", categoryID).Error; err != nil {
			return errorx.ErrDatabaseQuery.WithDetails(
				fmt.Sprintf("failed to find category with ID %v: %v", categoryID, err),
			)
		}

		if err := tx.WithContext(ctx).Model(entity).Association("Categories").Append(&category); err != nil {
			return errorx.ErrDatabaseQuery.WithDetails(
				fmt.Sprintf("failed to associate category %v with menu item: %v", categoryID, err),
			)
		}
	}
	return nil
}
