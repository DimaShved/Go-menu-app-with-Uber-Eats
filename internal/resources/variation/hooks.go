package variation

import (
	"context"

	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/platform/crud"
)

type Hooks struct {
	crud.NoopHooks[domain.Variation, CreateRequest, UpdateRequest, domain.Variation]
}

func (Hooks) AfterCreate(ctx context.Context, hookCtx crud.HookContext, request *CreateRequest, entity *domain.Variation) error {
	if len(request.Options) == 0 {
		return nil
	}

	options := make([]domain.VariationOption, 0, len(request.Options))
	for _, option := range request.Options {
		options = append(options, domain.VariationOption{
			Name:        option.Name,
			Price:       option.Price,
			IsAvailable: option.IsAvailable,
			VariationID: entity.ID,
		})
	}

	if err := hookCtx.Tx.WithContext(ctx).Create(&options).Error; err != nil {
		return errorx.ErrDatabaseQuery.WithDetails(err.Error())
	}
	entity.Options = options
	return nil
}
