package service

import (
	"uber-go-menu-copy/internal/domain"
	"uber-go-menu-copy/internal/repository"
)

type MenuSectionService struct {
	*genericService[*domain.MenuSections]
}

func NewMenuSectionService(repo repository.IGenericRepository[*domain.MenuSections]) *MenuSectionService {
	baseService := NewGenericService[*domain.MenuSections](repo)
	return &MenuSectionService{genericService: baseService.(*genericService[*domain.MenuSections])}
}
