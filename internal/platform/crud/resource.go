package crud

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
)

const DefaultTimeout = 5 * time.Second

type Validator interface {
	Struct(s any) error
}

type CreateMapper[Entity SoftDeleteEntity, CreateRequest any] func(CreateRequest) (*Entity, error)
type UpdateMapper[Entity SoftDeleteEntity, UpdateRequest any] func(*Entity, UpdateRequest) error
type ResponseMapper[Entity SoftDeleteEntity, Response any] func(*Entity) (Response, error)
type ExtraRoutes func(fiber.Router)

type Resource[Entity SoftDeleteEntity, CreateRequest any, UpdateRequest any, Response any] struct {
	Name        string
	Path        string
	Repository  Repository[Entity]
	TxManager   TxManager
	Validator   Validator
	Hooks       Hooks[Entity, CreateRequest, UpdateRequest, Response]
	MapCreate   CreateMapper[Entity, CreateRequest]
	ApplyUpdate UpdateMapper[Entity, UpdateRequest]
	MapResponse ResponseMapper[Entity, Response]
	ExtraRoutes ExtraRoutes
	Timeout     time.Duration
}

func (r Resource[Entity, CreateRequest, UpdateRequest, Response]) withDefaults() Resource[Entity, CreateRequest, UpdateRequest, Response] {
	if isNil(r.Hooks) {
		r.Hooks = NoopHooks[Entity, CreateRequest, UpdateRequest, Response]{}
	}
	if r.Timeout == 0 {
		r.Timeout = DefaultTimeout
	}
	return r
}

func (r Resource[Entity, CreateRequest, UpdateRequest, Response]) prepareForService() Resource[Entity, CreateRequest, UpdateRequest, Response] {
	r = r.withDefaults()
	r.mustValidate(false)
	return r
}

func (r Resource[Entity, CreateRequest, UpdateRequest, Response]) prepareForHandler() Resource[Entity, CreateRequest, UpdateRequest, Response] {
	r = r.withDefaults()
	r.mustValidate(true)
	return r
}

func (r Resource[Entity, CreateRequest, UpdateRequest, Response]) mustValidate(requirePath bool) {
	if err := r.validate(requirePath); err != nil {
		panic(err)
	}
}

func (r Resource[Entity, CreateRequest, UpdateRequest, Response]) validate(requirePath bool) error {
	missing := make([]string, 0)
	if strings.TrimSpace(r.Name) == "" {
		missing = append(missing, "Name")
	}
	if requirePath && strings.TrimSpace(r.Path) == "" {
		missing = append(missing, "Path")
	}
	if isNil(r.Repository) {
		missing = append(missing, "Repository")
	}
	if isNil(r.TxManager) {
		missing = append(missing, "TxManager")
	}
	if r.MapCreate == nil {
		missing = append(missing, "MapCreate")
	}
	if r.ApplyUpdate == nil {
		missing = append(missing, "ApplyUpdate")
	}
	if r.MapResponse == nil {
		missing = append(missing, "MapResponse")
	}

	if len(missing) == 0 {
		return nil
	}
	return fmt.Errorf("invalid CRUD resource %q: missing required field(s): %s", resourceName(r.Name), strings.Join(missing, ", "))
}

func resourceName(name string) string {
	if strings.TrimSpace(name) == "" {
		return "<unnamed>"
	}
	return name
}

func isNil(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
