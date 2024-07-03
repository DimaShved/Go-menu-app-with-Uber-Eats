package repository

import (
	"context"
	"github.com/google/uuid"
)

type IGenericRepository[T any] interface {
	Save(ctx context.Context, entity T, preloads ...string) error
	FindAll(ctx context.Context) ([]T, error)
	FindOneById(ctx context.Context, id uuid.UUID) (*T, error)
	HardDelete(ctx context.Context, id uuid.UUID) error

	// up to you, but may be useful
	FindAllIncludingDeleted(ctx context.Context) ([]T, error)
	FindOneByIdIncludingDeleted(ctx context.Context, id uuid.UUID) (*T, error)
	SoftDelete(ctx context.Context, id uuid.UUID) error
}
