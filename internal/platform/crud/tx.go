package crud

import (
	"context"

	"gorm.io/gorm"
)

type TxManager interface {
	DB() *gorm.DB
	WithinTx(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type GormTxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *GormTxManager {
	return &GormTxManager{db: db}
}

func (m *GormTxManager) DB() *gorm.DB {
	return m.db
}

func (m *GormTxManager) WithinTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}
