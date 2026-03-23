package modifier

import (
	"context"

	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.Modifier]
}

func NewRepository() *Repository {
	return &Repository{
		GormRepository: crud.NewGormRepository[domain.Modifier](crud.QueryOptions{
			Preloads: []string{"Options"},
		}),
	}
}

func (r *Repository) CreateOptions(ctx context.Context, tx *gorm.DB, options []domain.ModifierOption) error {
	if len(options) == 0 {
		return nil
	}
	if err := tx.WithContext(ctx).Create(&options).Error; err != nil {
		return errorx.ErrDatabaseQuery.WithDetails(err.Error())
	}
	return nil
}
