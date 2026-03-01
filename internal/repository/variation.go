package repository

import (
	"context"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
)

type VariationRepo interface {
	IGenericRepository[*domain.Variation]
	DB() *gorm.DB
	SaveTx(ctx context.Context, tx *gorm.DB, variation *domain.Variation) error
}

type variationRepo struct {
	*GenericRepository[*domain.Variation]
}

func NewVariationRepo(db *gorm.DB) VariationRepo {
	return &variationRepo{NewGenericRepo[*domain.Variation](db)}
}

func (r *variationRepo) DB() *gorm.DB {
	return r.db
}

func (r *variationRepo) SaveTx(ctx context.Context, tx *gorm.DB, variation *domain.Variation) error {
	return _handleDbError(tx.WithContext(ctx).Save(variation).Error, "SaveTx")
}
