package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log/slog"
	"time"
	"uber-go-menu-copy/internal/domain/interfaces"
)

type GenericRepository[T interfaces.Identifiable] struct {
	db *gorm.DB
}

func NewGenericRepo[T interfaces.Identifiable](db *gorm.DB) *GenericRepository[T] {
	return &GenericRepository[T]{db: db}
}

func (r *GenericRepository[T]) Save(ctx context.Context, entity T, preloads ...string) error {
	if _, ok := any(entity).(interfaces.Identifiable); !ok {
		slog.Error("entity does not implement Identifiable interface")
		return fmt.Errorf("entity does not implement Identifiable interface")
	}

	if err := r.db.WithContext(ctx).Save(&entity).Error; err != nil {
		return err
	}

	query := r.db.WithContext(ctx)
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	return query.First(&entity, "id = ?", entity.GetID()).Error
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
	if err != nil {
		return nil, err
	}
	return entities, nil
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
		return nil, err
	}
	return &entity, nil
}

func (r *GenericRepository[T]) HardDelete(ctx context.Context, id uuid.UUID) error {
	var model T
	return r.db.WithContext(ctx).Unscoped().Delete(&model, "id = ?", id).Error
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
	if err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *GenericRepository[T]) FindOneByIdIncludingDeleted(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx).Unscoped()

	if preloader, ok := any(entity).(interfaces.Preloader); ok {
		for _, relation := range preloader.PreloadRelations() {
			query = query.Preload(relation)
		}
	}
	err := query.Find(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GenericRepository[T]) SoftDelete(ctx context.Context, id uuid.UUID) error {
	var model T
	return r.db.WithContext(ctx).Model(&model).Where("id = ?", id).Update("deleted_at", time.Now()).Error
}
