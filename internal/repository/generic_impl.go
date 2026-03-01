package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
	"uber-go-menu/internal/domain/interfaces"
	"uber-go-menu/internal/pkg/errorx"
)

type GenericRepository[T interfaces.Identifiable] struct {
	db *gorm.DB
}

func NewGenericRepo[T interfaces.Identifiable](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

func _handleDbError(err error, operation string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorx.ErrRecordNotFound.Wrap(err)
	}

	wrappedErr := fmt.Errorf("operation '%s' failed: %w", operation, err)
	return errorx.ErrDatabaseQuery.Wrap(wrappedErr)
}

func (r *GenericRepository[T]) Save(ctx context.Context, entity *T) error {
	return _handleDbError(r.db.WithContext(ctx).Save(entity).Error, "Save")
}

func (r *GenericRepository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	var entity T
	query := r.db.WithContext(ctx).Model(entity)

	if preloader, ok := any(entity).(interfaces.Preloader); ok {
		for _, relation := range preloader.PreloadRelations() {
			query = query.Preload(relation)
		}
	}
	err := query.Where("deleted_at IS NULL").Find(&entities).Error
	return entities, _handleDbError(err, "FindAll")
}

func (r *GenericRepository[T]) FindAllIncludingDeleted(ctx context.Context) ([]T, error) {
	var entities []T
	var entity T
	query := r.db.WithContext(ctx).Unscoped().Model(entity)

	if preloader, ok := any(entity).(interfaces.Preloader); ok {
		for _, relation := range preloader.PreloadRelations() {
			query = query.Preload(relation)
		}
	}

	err := query.Find(&entities).Error
	return entities, _handleDbError(err, "FindAllIncludingDeleted")
}

func (r *GenericRepository[T]) FindOneById(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx).Where("deleted_at IS NULL")

	if preloader, ok := any(entity).(interfaces.Preloader); ok {
		for _, relation := range preloader.PreloadRelations() {
			query = query.Preload(relation)
		}
	}

	err := query.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, _handleDbError(err, "FindOneById")
	}
	return &entity, nil
}

func (r *GenericRepository[T]) FindOneByIdIncludingDeleted(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx).Unscoped()

	if preloader, ok := any(entity).(interfaces.Preloader); ok {
		for _, relation := range preloader.PreloadRelations() {
			query = query.Preload(relation)
		}
	}
	err := query.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, _handleDbError(err, "FindOneByIdIncludingDeleted")
	}
	return &entity, nil
}

func (r *GenericRepository[T]) SoftDelete(ctx context.Context, id uuid.UUID) error {
	var model T
	result := r.db.WithContext(ctx).Model(&model).Where("id = ? AND deleted_at IS NULL", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		return _handleDbError(result.Error, "SoftDelete")
	}
	if result.RowsAffected == 0 {
		return errorx.ErrRecordNotFound.Wrap(errors.New("record not found or already deleted"))
	}
	return nil
}

func (r *GenericRepository[T]) HardDelete(ctx context.Context, id uuid.UUID) error {
	var model T
	result := r.db.WithContext(ctx).Unscoped().Delete(&model, "id = ?", id)
	if result.Error != nil {
		return _handleDbError(result.Error, "HardDelete")
	}
	return nil
}
