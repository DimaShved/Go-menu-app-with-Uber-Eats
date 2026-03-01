package repository

import (
	"context"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
)

type ModifierRepo interface {
	IGenericRepository[*domain.Modifier]
	DB() *gorm.DB
	SaveTx(ctx context.Context, tx *gorm.DB, modifier *domain.Modifier) error
}

type modifierRepo struct {
	*GenericRepository[*domain.Modifier]
}

func NewModifierRepo(db *gorm.DB) ModifierRepo {
	return &modifierRepo{NewGenericRepo[*domain.Modifier](db)}
}

func (r *modifierRepo) DB() *gorm.DB {
	return r.db
}

func (r *modifierRepo) SaveTx(ctx context.Context, tx *gorm.DB, modifier *domain.Modifier) error {
	return _handleDbError(tx.WithContext(ctx).Save(modifier).Error, "SaveTx")
}
