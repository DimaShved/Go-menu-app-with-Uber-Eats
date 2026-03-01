package service

import (
	"context"
	"gorm.io/gorm"
	"log"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/inputs"
	"uber-go-menu/internal/repository"
)

type ModifierService struct {
	genericService        IGenericService[*domain.Modifier]
	modifierRepo          repository.ModifierRepo
	modifierOptionService ModifierOptionService
}

func (ms *ModifierService) GenericService() IGenericService[*domain.Modifier] {
	return ms.genericService
}

func NewModifierService(repo repository.ModifierRepo, vos *ModifierOptionService) *ModifierService {
	genericService := NewGenericService[*domain.Modifier](repo)
	return &ModifierService{
		genericService:        genericService,
		modifierRepo:          repo,
		modifierOptionService: *vos,
	}
}

func (ms *ModifierService) CreateWithOptions(ctx context.Context, modifierInput *inputs.ModifierInput) (*domain.Modifier, error) {
	log.Printf("Received modifier input: %+v", modifierInput)

	modifier := &domain.Modifier{
		Name:              modifierInput.Name,
		TotalMaxSelection: modifierInput.TotalMaxSelection,
		CategoryID:        modifierInput.CategoryID,
	}

	err := ms.modifierRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := ms.modifierRepo.SaveTx(ctx, tx, modifier); err != nil {
			return err
		}

		if len(modifierInput.Options) == 0 {
			return nil
		}

		options := make([]domain.ModifierOption, 0, len(modifierInput.Options))
		for _, inputOption := range modifierInput.Options {
			options = append(options, domain.ModifierOption{
				Name:         inputOption.Name,
				Price:        inputOption.Price,
				MaxSelection: inputOption.MaxSelection,
				IsAvailable:  inputOption.IsAvailable,
				ModifierID:   modifier.ID,
			})
		}

		savedOptions, err := ms.modifierOptionService.CreateManyTx(ctx, tx, options)
		if err != nil {
			return err
		}

		modifier.Options = savedOptions
		return nil
	})
	if err != nil {
		return nil, err
	}

	return modifier, nil
}
