package crud

import (
	"time"

	"github.com/google/uuid"
)

const DefaultTimeout = 5 * time.Second

type Validator interface {
	Struct(s any) error
}

type IDFunc[Entity any] func(*Entity) uuid.UUID
type CreateMapper[Entity any, CreateRequest any] func(CreateRequest) (*Entity, error)
type UpdateMapper[Entity any, UpdateRequest any] func(*Entity, UpdateRequest) error
type ResponseMapper[Entity any, Response any] func(*Entity) (Response, error)

type Resource[Entity any, CreateRequest any, UpdateRequest any, Response any] struct {
	Name        string
	Path        string
	Repository  Repository[Entity]
	TxManager   TxManager
	Validator   Validator
	Hooks       Hooks[Entity, CreateRequest, UpdateRequest, Response]
	GetID       IDFunc[Entity]
	MapCreate   CreateMapper[Entity, CreateRequest]
	ApplyUpdate UpdateMapper[Entity, UpdateRequest]
	MapResponse ResponseMapper[Entity, Response]
	Timeout     time.Duration
}

func (r Resource[Entity, CreateRequest, UpdateRequest, Response]) withDefaults() Resource[Entity, CreateRequest, UpdateRequest, Response] {
	if r.Hooks == nil {
		r.Hooks = NoopHooks[Entity, CreateRequest, UpdateRequest, Response]{}
	}
	if r.Timeout == 0 {
		r.Timeout = DefaultTimeout
	}
	return r
}
