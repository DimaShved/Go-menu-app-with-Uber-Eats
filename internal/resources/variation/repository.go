package variation

import (
	"context"

	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.Variation]
}

func NewRepository() *Repository {
	return &Repository{
		GormRepository: crud.NewGormRepository[domain.Variation](crud.QueryOptions{
			Preloads: []string{"Options", "Category.Section", "Category.Section.Restaurant"},
		}),
	}
}

func (r *Repository) CreateOptions(ctx context.Context, tx *gorm.DB, options []domain.VariationOption) error {
	if len(options) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).Create(&options).Error; err != nil {
		return errorx.ErrDatabaseQuery.WithDetails(err.Error())
	}
	return nil
}
