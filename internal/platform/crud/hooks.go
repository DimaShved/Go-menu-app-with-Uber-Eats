package crud

import (
	"context"

	"gorm.io/gorm"
)

type HookContext struct {
	ResourceName string
	Tx           *gorm.DB
}

type Hooks[Entity any, CreateRequest any, UpdateRequest any, Response any] interface {
	BeforeCreate(context.Context, HookContext, *CreateRequest, *Entity) error
	AfterCreate(context.Context, HookContext, *CreateRequest, *Entity) error
	BeforeUpdate(context.Context, HookContext, *UpdateRequest, *Entity) error
	AfterUpdate(context.Context, HookContext, *UpdateRequest, *Entity) error
	BeforeDelete(context.Context, HookContext, *Entity) error
	AfterDelete(context.Context, HookContext, *Entity) error
	BeforeResponse(context.Context, HookContext, *Entity, *Response) error
}

type NoopHooks[Entity any, CreateRequest any, UpdateRequest any, Response any] struct{}

func (NoopHooks[Entity, CreateRequest, UpdateRequest, Response]) BeforeCreate(context.Context, HookContext, *CreateRequest, *Entity) error {
	return nil
}

func (NoopHooks[Entity, CreateRequest, UpdateRequest, Response]) AfterCreate(context.Context, HookContext, *CreateRequest, *Entity) error {
	return nil
}

func (NoopHooks[Entity, CreateRequest, UpdateRequest, Response]) BeforeUpdate(context.Context, HookContext, *UpdateRequest, *Entity) error {
	return nil
}

func (NoopHooks[Entity, CreateRequest, UpdateRequest, Response]) AfterUpdate(context.Context, HookContext, *UpdateRequest, *Entity) error {
	return nil
}

func (NoopHooks[Entity, CreateRequest, UpdateRequest, Response]) BeforeDelete(context.Context, HookContext, *Entity) error {
	return nil
}

func (NoopHooks[Entity, CreateRequest, UpdateRequest, Response]) AfterDelete(context.Context, HookContext, *Entity) error {
	return nil
}

func (NoopHooks[Entity, CreateRequest, UpdateRequest, Response]) BeforeResponse(context.Context, HookContext, *Entity, *Response) error {
	return nil
}
