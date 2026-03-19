package crud

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"uber-go-menu/internal/pkg/errorx"
)

type ResourceService[Entity any, CreateRequest any, UpdateRequest any, Response any] struct {
	resource Resource[Entity, CreateRequest, UpdateRequest, Response]
}

func NewService[Entity any, CreateRequest any, UpdateRequest any, Response any](resource Resource[Entity, CreateRequest, UpdateRequest, Response]) *ResourceService[Entity, CreateRequest, UpdateRequest, Response] {
	return &ResourceService[Entity, CreateRequest, UpdateRequest, Response]{
		resource: resource.withDefaults(),
	}
}

func (s *ResourceService[Entity, CreateRequest, UpdateRequest, Response]) Create(ctx context.Context, request CreateRequest) (Response, error) {
	var zero Response
	if err := s.validate(request); err != nil {
		return zero, err
	}

	entity, err := s.resource.MapCreate(request)
	if err != nil {
		return zero, err
	}

	var created *Entity
	if err := s.resource.TxManager.WithinTx(ctx, func(tx *gorm.DB) error {
		hookCtx := HookContext{ResourceName: s.resource.Name, Tx: tx}
		if err := s.resource.Hooks.BeforeCreate(ctx, hookCtx, &request, entity); err != nil {
			return err
		}
		if err := s.resource.Repository.Create(ctx, tx, entity); err != nil {
			return err
		}
		if err := s.resource.Hooks.AfterCreate(ctx, hookCtx, &request, entity); err != nil {
			return err
		}

		reloaded, err := s.resource.Repository.GetByID(ctx, tx, s.resource.GetID(entity))
		if err != nil {
			return err
		}
		created = reloaded
		return nil
	}); err != nil {
		return zero, err
	}

	return s.response(ctx, created)
}

func (s *ResourceService[Entity, CreateRequest, UpdateRequest, Response]) Update(ctx context.Context, id uuid.UUID, request UpdateRequest) (Response, error) {
	var zero Response
	if err := s.validate(request); err != nil {
		return zero, err
	}

	var updated *Entity
	if err := s.resource.TxManager.WithinTx(ctx, func(tx *gorm.DB) error {
		entity, err := s.resource.Repository.GetByID(ctx, tx, id)
		if err != nil {
			return err
		}
		if err := s.resource.ApplyUpdate(entity, request); err != nil {
			return err
		}

		hookCtx := HookContext{ResourceName: s.resource.Name, Tx: tx}
		if err := s.resource.Hooks.BeforeUpdate(ctx, hookCtx, &request, entity); err != nil {
			return err
		}
		if err := s.resource.Repository.Update(ctx, tx, entity); err != nil {
			return err
		}
		if err := s.resource.Hooks.AfterUpdate(ctx, hookCtx, &request, entity); err != nil {
			return err
		}

		reloaded, err := s.resource.Repository.GetByID(ctx, tx, id)
		if err != nil {
			return err
		}
		updated = reloaded
		return nil
	}); err != nil {
		return zero, err
	}

	return s.response(ctx, updated)
}

func (s *ResourceService[Entity, CreateRequest, UpdateRequest, Response]) Delete(ctx context.Context, id uuid.UUID) error {
	return s.resource.TxManager.WithinTx(ctx, func(tx *gorm.DB) error {
		entity, err := s.resource.Repository.GetByID(ctx, tx, id)
		if err != nil {
			return err
		}

		hookCtx := HookContext{ResourceName: s.resource.Name, Tx: tx}
		if err := s.resource.Hooks.BeforeDelete(ctx, hookCtx, entity); err != nil {
			return err
		}
		if err := s.resource.Repository.Delete(ctx, tx, id); err != nil {
			return err
		}
		if err := s.resource.Hooks.AfterDelete(ctx, hookCtx, entity); err != nil {
			return err
		}
		return nil
	})
}

func (s *ResourceService[Entity, CreateRequest, UpdateRequest, Response]) GetByID(ctx context.Context, id uuid.UUID) (Response, error) {
	var zero Response
	entity, err := s.resource.Repository.GetByID(ctx, s.resource.TxManager.DB(), id)
	if err != nil {
		return zero, err
	}
	return s.response(ctx, entity)
}

func (s *ResourceService[Entity, CreateRequest, UpdateRequest, Response]) List(ctx context.Context) ([]Response, error) {
	entities, err := s.resource.Repository.List(ctx, s.resource.TxManager.DB())
	if err != nil {
		return nil, err
	}

	responses := make([]Response, 0, len(entities))
	for i := range entities {
		response, err := s.response(ctx, &entities[i])
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (s *ResourceService[Entity, CreateRequest, UpdateRequest, Response]) validate(request any) error {
	if s.resource.Validator == nil {
		return nil
	}
	if err := s.resource.Validator.Struct(request); err != nil {
		return errorx.ErrInvalidInput.WithDetails(fmt.Sprintf("Input validation failed: %v", err))
	}
	return nil
}

func (s *ResourceService[Entity, CreateRequest, UpdateRequest, Response]) response(ctx context.Context, entity *Entity) (Response, error) {
	response, err := s.resource.MapResponse(entity)
	if err != nil {
		var zero Response
		return zero, err
	}
	hookCtx := HookContext{ResourceName: s.resource.Name}
	if err := s.resource.Hooks.BeforeResponse(ctx, hookCtx, entity, &response); err != nil {
		var zero Response
		return zero, err
	}
	return response, nil
}
