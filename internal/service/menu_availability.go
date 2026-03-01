package service

import (
	"fmt"
	"log"
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/pkg/errorx"
	"uber-go-menu/internal/repository"
)

type MenuAvailabilityService struct {
	genericService       IGenericService[*domain.MenuAvailability]
	menuAvailabilityRepo repository.MenuAvailabilityRepo
}

func (s *MenuAvailabilityService) GenericService() IGenericService[*domain.MenuAvailability] {
	return s.genericService
}

func NewMenuAvailabilityService(menuAvailabilityRepo repository.MenuAvailabilityRepo) *MenuAvailabilityService {
	genericService := NewGenericService[*domain.MenuAvailability](menuAvailabilityRepo)
	return &MenuAvailabilityService{
		genericService:       genericService,
		menuAvailabilityRepo: menuAvailabilityRepo,
	}
}

func (s *MenuAvailabilityService) UpsertBatch(menuAvailabilities []*domain.MenuAvailability) ([]*domain.MenuAvailability, error) {
	existingEntries := make(map[string]bool)
	for _, ma := range menuAvailabilities {
		key := ma.MenuSectionId.String() + fmt.Sprintf("-%d", ma.DayOfWeek)
		if existingEntries[key] {
			log.Printf("Duplicate entry for restaurant %s on day %d", ma.MenuSectionId, ma.DayOfWeek)
			return nil, errorx.NewAppError(106, fmt.Sprintf("duplicate entry for restaurant %s on day %d", ma.MenuSectionId, ma.DayOfWeek))
		}
		existingEntries[key] = true
	}

	updatedAvailabilities, err := s.menuAvailabilityRepo.UpsertAvailabilities(menuAvailabilities)
	if err != nil {
		return nil, err
	}
	return updatedAvailabilities, nil
}
