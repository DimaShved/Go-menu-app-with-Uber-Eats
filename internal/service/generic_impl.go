package service

import (
	"context"
	"github.com/google/uuid"
	"uber-go-menu/internal/domain/interfaces"
	"uber-go-menu/internal/repository"
)

type genericService[T interfaces.Identifiable] struct {
	repo repository.IGenericRepository[T]
}

func NewGenericService[T interfaces.Identifiable](repo repository.IGenericRepository[T]) IGenericService[T] {
	return &genericService[T]{repo: repo}
}

func (s *genericService[T]) Save(ctx context.Context, entity T) (T, error) {
	err := s.repo.Save(ctx, &entity)
	if err != nil {
		var zero T
		return zero, err
	}

	fetchedEntity, err := s.repo.FindOneById(ctx, entity.GetID())
	if err != nil {
		var zero T
		return zero, err
	}

	return *fetchedEntity, nil
}

func (s *genericService[T]) FindAll(ctx context.Context) ([]T, error) {
	return s.repo.FindAll(ctx)
}

func (s *genericService[T]) FindAllIncludingDeleted(ctx context.Context) ([]T, error) {
	return s.repo.FindAllIncludingDeleted(ctx)
}

func (s *genericService[T]) FindOneById(ctx context.Context, id uuid.UUID) (*T, error) {
	return s.repo.FindOneById(ctx, id)
}

func (s *genericService[T]) FindOneByIdIncludingDeleted(ctx context.Context, id uuid.UUID) (*T, error) {
	return s.repo.FindOneByIdIncludingDeleted(ctx, id)
}

func (s *genericService[T]) SoftDelete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}

func (s *genericService[T]) HardDelete(ctx context.Context, id uuid.UUID) error {
	return s.repo.HardDelete(ctx, id)
}
