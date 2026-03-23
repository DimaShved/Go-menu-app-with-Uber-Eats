package menu_item

import (
	"context"

	"github.com/google/uuid"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

type Hooks struct {
	crud.NoopHooks[domain.MenuItem, CreateRequest, UpdateRequest, Response]
	repository *Repository
}

func (h Hooks) AfterCreate(ctx context.Context, hookCtx crud.HookContext, request *CreateRequest, entity *domain.MenuItem) error {
	if len(request.Categories) == 0 {
		return nil
	}

	categoryIDs := make([]uuid.UUID, 0, len(request.Categories))
	for _, id := range request.Categories {
		categoryID, err := uuid.Parse(id)
		if err != nil {
			return errorx.ErrInvalidInput.WithDetails("invalid category ID: " + id)
		}
		categoryIDs = append(categoryIDs, categoryID)
	}

	return h.repository.AttachCategories(ctx, hookCtx.Tx, entity, categoryIDs)
}
