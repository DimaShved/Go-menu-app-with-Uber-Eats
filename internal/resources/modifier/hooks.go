package modifier

import (
	"context"

	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

type Hooks struct {
	crud.NoopHooks[domain.Modifier, CreateRequest, UpdateRequest, Response]
	repository *Repository
}

func (h Hooks) AfterCreate(ctx context.Context, hookCtx crud.HookContext, request *CreateRequest, entity *domain.Modifier) error {
	if len(request.Options) == 0 {
		return nil
	}

	options := make([]domain.ModifierOption, 0, len(request.Options))
	for _, option := range request.Options {
		options = append(options, domain.ModifierOption{
			Name:         option.Name,
			Price:        option.Price,
			MaxSelection: option.MaxSelection,
			IsAvailable:  option.IsAvailable,
			ModifierID:   entity.ID,
		})
	}

	if err := h.repository.CreateOptions(ctx, hookCtx.Tx, options); err != nil {
		return err
	}
	entity.Options = options
	return nil
}
