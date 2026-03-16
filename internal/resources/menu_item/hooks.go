package menu_item

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

type Hooks struct {
	crud.NoopHooks[domain.MenuItem, CreateRequest, UpdateRequest, domain.MenuItem]
}

func (Hooks) AfterCreate(ctx context.Context, hookCtx crud.HookContext, request *CreateRequest, entity *domain.MenuItem) error {
	for _, id := range request.Categories {
		categoryID, err := uuid.Parse(id)
		if err != nil {
			return errorx.ErrInvalidInput.WithDetails("invalid category ID: " + id)
		}

		var category domain.MenuCategory
		if err := hookCtx.Tx.WithContext(ctx).First(&category, "id = ?", categoryID).Error; err != nil {
			return errorx.ErrDatabaseQuery.WithDetails(
				fmt.Sprintf("failed to find category with ID %v: %v", categoryID, err),
			)
		}

		if err := hookCtx.Tx.WithContext(ctx).Model(entity).Association("Categories").Append(&category); err != nil {
			return errorx.ErrDatabaseQuery.WithDetails(
				fmt.Sprintf("failed to associate category %v with menu item: %v", categoryID, err),
			)
		}
	}
	return nil
}
