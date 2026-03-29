package modifier_option

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/google/uuid"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

func TestNewResourceMapsCreateRequest(t *testing.T) {
	resource := resourceFromNew(t)
	modifierID := uuid.New()
	request := CreateRequest{
		Name:         "Extra cheese",
		Price:        250,
		MaxSelection: 2,
		IsAvailable:  true,
		ModifierID:   modifierID,
	}

	entity, err := resource.MapCreate(request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if entity.Name != request.Name {
		t.Fatalf("expected Name %q, got %q", request.Name, entity.Name)
	}
	if entity.Price != request.Price {
		t.Fatalf("expected Price %d, got %d", request.Price, entity.Price)
	}
	if entity.MaxSelection != request.MaxSelection {
		t.Fatalf("expected MaxSelection %d, got %d", request.MaxSelection, entity.MaxSelection)
	}
	if entity.IsAvailable != request.IsAvailable {
		t.Fatalf("expected IsAvailable %v, got %v", request.IsAvailable, entity.IsAvailable)
	}
	if entity.ModifierID != request.ModifierID {
		t.Fatalf("expected ModifierID %s, got %s", request.ModifierID, entity.ModifierID)
	}
}

func TestNewResourceAppliesUpdateRequest(t *testing.T) {
	resource := resourceFromNew(t)
	modifierID := uuid.New()
	entity := &domain.ModifierOption{Name: "Old option"}
	request := UpdateRequest{
		Name:         "No onions",
		Price:        0,
		MaxSelection: 1,
		IsAvailable:  true,
		ModifierID:   modifierID,
	}

	err := resource.ApplyUpdate(entity, request)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if entity.Name != request.Name {
		t.Fatalf("expected Name %q, got %q", request.Name, entity.Name)
	}
	if entity.Price != request.Price {
		t.Fatalf("expected Price %d, got %d", request.Price, entity.Price)
	}
	if entity.MaxSelection != request.MaxSelection {
		t.Fatalf("expected MaxSelection %d, got %d", request.MaxSelection, entity.MaxSelection)
	}
	if entity.IsAvailable != request.IsAvailable {
		t.Fatalf("expected IsAvailable %v, got %v", request.IsAvailable, entity.IsAvailable)
	}
	if entity.ModifierID != request.ModifierID {
		t.Fatalf("expected ModifierID %s, got %s", request.ModifierID, entity.ModifierID)
	}
}

func resourceFromNew(t *testing.T) crud.Resource[domain.ModifierOption, CreateRequest, UpdateRequest, Response] {
	t.Helper()

	handler, ok := New(nil, nil).(*crud.ResourceHandler[domain.ModifierOption, CreateRequest, UpdateRequest, Response])
	if !ok {
		t.Fatalf("expected CRUD resource handler, got %T", New(nil, nil))
	}

	field := reflect.ValueOf(handler).Elem().FieldByName("resource")
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Interface().(crud.Resource[domain.ModifierOption, CreateRequest, UpdateRequest, Response])
}
