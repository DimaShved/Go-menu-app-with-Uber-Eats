package repository

import (
	"context"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
)

type ModifierOptionRepo interface {
	IGenericRepository[*domain.ModifierOption]
	CreateMany(options []domain.ModifierOption) ([]domain.ModifierOption, error)
	CreateManyTx(ctx context.Context, tx *gorm.DB, options []domain.ModifierOption) ([]domain.ModifierOption, error)
}

type modifierOptionRepo struct {
	*GenericRepository[*domain.ModifierOption]
}

func NewModifierOptionRepo(db *gorm.DB) ModifierOptionRepo {
	return &modifierOptionRepo{NewGenericRepo[*domain.ModifierOption](db)}
}

func (r *modifierOptionRepo) CreateMany(options []domain.ModifierOption) ([]domain.ModifierOption, error) {
	return r.CreateManyTx(context.Background(), r.db, options)
}

func (r *modifierOptionRepo) CreateManyTx(ctx context.Context, tx *gorm.DB, options []domain.ModifierOption) ([]domain.ModifierOption, error) {
	if err := tx.WithContext(ctx).Create(&options).Error; err != nil {
		return nil, _handleDbError(err, "CreateManyTx")
	}
	return options, nil
}
