package service

import (
	"context"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/repository"
)

type ModifierOptionService struct {
	genericService     IGenericService[*domain.ModifierOption]
	modifierOptionRepo repository.ModifierOptionRepo
}

func (vos *ModifierOptionService) GenericService() IGenericService[*domain.ModifierOption] {
	return vos.genericService
}

func NewModifierOptionService(repo repository.ModifierOptionRepo) *ModifierOptionService {
	genericService := NewGenericService[*domain.ModifierOption](repo)
	return &ModifierOptionService{
		genericService:     genericService,
		modifierOptionRepo: repo,
	}
}

func (vos *ModifierOptionService) CreateMany(options []domain.ModifierOption) ([]domain.ModifierOption, error) {
	savedOptions, err := vos.modifierOptionRepo.CreateMany(options)
	if err != nil {
		return nil, err
	}
	return savedOptions, nil
}

func (vos *ModifierOptionService) CreateManyTx(ctx context.Context, tx *gorm.DB, options []domain.ModifierOption) ([]domain.ModifierOption, error) {
	return vos.modifierOptionRepo.CreateManyTx(ctx, tx, options)
}
