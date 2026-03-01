package repository

import (
	"context"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
)

type VariationOptionRepo interface {
	IGenericRepository[*domain.VariationOption]
	CreateMany(options []domain.VariationOption) ([]domain.VariationOption, error)
	CreateManyTx(ctx context.Context, tx *gorm.DB, options []domain.VariationOption) ([]domain.VariationOption, error)
}

type variationOptionRepo struct {
	*GenericRepository[*domain.VariationOption]
}

func NewVariationOptionRepo(db *gorm.DB) VariationOptionRepo {
	return &variationOptionRepo{NewGenericRepo[*domain.VariationOption](db)}
}

func (r *variationOptionRepo) CreateMany(options []domain.VariationOption) ([]domain.VariationOption, error) {
	return r.CreateManyTx(context.Background(), r.db, options)
}

func (r *variationOptionRepo) CreateManyTx(ctx context.Context, tx *gorm.DB, options []domain.VariationOption) ([]domain.VariationOption, error) {
	if err := tx.WithContext(ctx).Create(&options).Error; err != nil {
		return nil, _handleDbError(err, "CreateManyTx")
	}
	return options, nil
}
