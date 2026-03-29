package menu_item

import (
	"context"
	"strings"
	"testing"

	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

func TestHooksAfterCreateSkipsAttachWhenNoCategoriesRequested(t *testing.T) {
	hooks := Hooks{}
	request := CreateRequest{}

	err := hooks.AfterCreate(context.Background(), crud.HookContext{}, &request, &domain.MenuItem{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestHooksAfterCreateRejectsInvalidCategoryID(t *testing.T) {
	hooks := Hooks{}
	request := CreateRequest{
		Categories: []string{"not-a-category-id"},
	}

	err := hooks.AfterCreate(context.Background(), crud.HookContext{}, &request, &domain.MenuItem{})
	if err == nil {
		t.Fatal("expected invalid input error")
	}

	apiErr, ok := err.(*errorx.APIError)
	if !ok || apiErr.Code != errorx.ErrInvalidInput.Code {
		t.Fatalf("expected invalid input error, got %T: %v", err, err)
	}
	if !strings.Contains(apiErr.Details, "not-a-category-id") {
		t.Fatalf("expected invalid category ID in details, got %q", apiErr.Details)
	}
}
