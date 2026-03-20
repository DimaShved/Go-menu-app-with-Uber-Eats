package crud

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/pkg/errorx"
)

type Repository[Entity any] interface {
	Create(ctx context.Context, tx *gorm.DB, entity *Entity) error
	Update(ctx context.Context, tx *gorm.DB, entity *Entity) error
	Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error
	GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*Entity, error)
	List(ctx context.Context, db *gorm.DB) ([]Entity, error)
}

type QueryOptions struct {
	Preloads []string
}

type GormRepository[Entity any] struct {
	queryOptions QueryOptions
}

func NewGormRepository[Entity any](options ...QueryOptions) *GormRepository[Entity] {
	var queryOptions QueryOptions
	if len(options) > 0 {
		queryOptions = options[0]
	}
	return &GormRepository[Entity]{queryOptions: queryOptions}
}

func (r *GormRepository[Entity]) Create(ctx context.Context, tx *gorm.DB, entity *Entity) error {
	return handleDBError(tx.WithContext(ctx).Create(entity).Error, "Create")
}

func (r *GormRepository[Entity]) Update(ctx context.Context, tx *gorm.DB, entity *Entity) error {
	return handleDBError(tx.WithContext(ctx).Save(entity).Error, "Update")
}

func (r *GormRepository[Entity]) Delete(ctx context.Context, tx *gorm.DB, id uuid.UUID) error {
	var entity Entity
	result := tx.WithContext(ctx).
		Model(&entity).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now())
	if result.Error != nil {
		return handleDBError(result.Error, "Delete")
	}
	if result.RowsAffected == 0 {
		return errorx.ErrRecordNotFound.Wrap(errors.New("record not found or already deleted"))
	}
	return nil
}

func (r *GormRepository[Entity]) GetByID(ctx context.Context, db *gorm.DB, id uuid.UUID) (*Entity, error) {
	var entity Entity
	query := r.applyQueryOptions(db.WithContext(ctx)).Where("deleted_at IS NULL")
	if err := query.First(&entity, "id = ?", id).Error; err != nil {
		return nil, handleDBError(err, "GetByID")
	}
	return &entity, nil
}

func (r *GormRepository[Entity]) List(ctx context.Context, db *gorm.DB) ([]Entity, error) {
	var entities []Entity
	query := r.applyQueryOptions(db.WithContext(ctx).Model(new(Entity))).Where("deleted_at IS NULL")
	err := query.Find(&entities).Error
	return entities, handleDBError(err, "List")
}

func (r *GormRepository[Entity]) applyQueryOptions(query *gorm.DB) *gorm.DB {
	for _, relation := range r.queryOptions.Preloads {
		query = query.Preload(relation)
	}
	return query
}

func handleDBError(err error, operation string) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errorx.ErrRecordNotFound.Wrap(err)
	}
	return errorx.ErrDatabaseQuery.Wrap(fmt.Errorf("operation %q failed: %w", operation, err))
}
