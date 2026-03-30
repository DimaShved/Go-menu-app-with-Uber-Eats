package variation_option

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
	variationID := uuid.New()
	request := CreateRequest{
		Name:        "Large",
		Price:       300,
		IsAvailable: true,
		VariationID: variationID,
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
	if entity.IsAvailable != request.IsAvailable {
		t.Fatalf("expected IsAvailable %v, got %v", request.IsAvailable, entity.IsAvailable)
	}
	if entity.VariationID != request.VariationID {
		t.Fatalf("expected VariationID %s, got %s", request.VariationID, entity.VariationID)
	}
}

func TestNewResourceAppliesUpdateRequest(t *testing.T) {
	resource := resourceFromNew(t)
	variationID := uuid.New()
	entity := &domain.VariationOption{Name: "Small"}
	request := UpdateRequest{
		Name:        "Medium",
		Price:       150,
		IsAvailable: false,
		VariationID: variationID,
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
	if entity.IsAvailable != request.IsAvailable {
		t.Fatalf("expected IsAvailable %v, got %v", request.IsAvailable, entity.IsAvailable)
	}
	if entity.VariationID != request.VariationID {
		t.Fatalf("expected VariationID %s, got %s", request.VariationID, entity.VariationID)
	}
}

func resourceFromNew(t *testing.T) crud.Resource[domain.VariationOption, CreateRequest, UpdateRequest, Response] {
	t.Helper()

	handler, ok := New(nil, nil).(*crud.ResourceHandler[domain.VariationOption, CreateRequest, UpdateRequest, Response])
	if !ok {
		t.Fatalf("expected CRUD resource handler, got %T", New(nil, nil))
	}

	field := reflect.ValueOf(handler).Elem().FieldByName("resource")
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).
		Elem().
		Interface().(crud.Resource[domain.VariationOption, CreateRequest, UpdateRequest, Response])
}
