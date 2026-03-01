package service

import (
	"context"
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/inputs"
	"uber-go-menu/internal/repository"
)

type VariationService struct {
	genericService         IGenericService[*domain.Variation]
	variationRepo          repository.VariationRepo
	variationOptionService VariationOptionService
}

func (s *VariationService) GenericService() IGenericService[*domain.Variation] {
	return s.genericService
}

func NewVariationService(repo repository.VariationRepo, vos *VariationOptionService) *VariationService {
	genericService := NewGenericService[*domain.Variation](repo)
	return &VariationService{
		genericService:         genericService,
		variationRepo:          repo,
		variationOptionService: *vos,
	}
}

func (s *VariationService) CreateWithOptions(ctx context.Context, variationInput *inputs.VariationInput) (*domain.Variation, error) {
	variation := &domain.Variation{
		Name:       variationInput.Name,
		CategoryID: variationInput.CategoryID,
	}

	err := s.variationRepo.DB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.variationRepo.SaveTx(ctx, tx, variation); err != nil {
			return err
		}

		if len(variationInput.Options) == 0 {
			return nil
		}

		options := make([]domain.VariationOption, 0, len(variationInput.Options))
		for _, inputOption := range variationInput.Options {
			options = append(options, domain.VariationOption{
				Name:        inputOption.Name,
				Price:       inputOption.Price,
				IsAvailable: inputOption.IsAvailable,
				VariationID: variation.ID,
			})
		}

		savedOptions, err := s.variationOptionService.CreateManyTx(ctx, tx, options)
		if err != nil {
			return err
		}

		variation.Options = savedOptions
		return nil
	})
	if err != nil {
		return nil, err
	}

	return variation, nil
}
