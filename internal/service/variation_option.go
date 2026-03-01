package service

import (
	"context"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/repository"
)

type VariationOptionService struct {
	genericService      IGenericService[*domain.VariationOption]
	variationOptionRepo repository.VariationOptionRepo
}

func (vos *VariationOptionService) GenericService() IGenericService[*domain.VariationOption] {
	return vos.genericService
}

func NewVariationOptionService(repo repository.VariationOptionRepo) *VariationOptionService {
	genericService := NewGenericService[*domain.VariationOption](repo)
	return &VariationOptionService{
		genericService:      genericService,
		variationOptionRepo: repo,
	}
}

func (vos *VariationOptionService) CreateMany(options []domain.VariationOption) ([]domain.VariationOption, error) {
	savedOptions, err := vos.variationOptionRepo.CreateMany(options)
	if err != nil {
		return nil, err
	}
	return savedOptions, nil
}

func (vos *VariationOptionService) CreateManyTx(ctx context.Context, tx *gorm.DB, options []domain.VariationOption) ([]domain.VariationOption, error) {
	return vos.variationOptionRepo.CreateManyTx(ctx, tx, options)
}
