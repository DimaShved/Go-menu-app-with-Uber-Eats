package service

import (
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/repository"
)

type MenuSectionService struct {
	*genericService[*domain.MenuSection]
}

func NewMenuSectionService(repo repository.IGenericRepository[*domain.MenuSection]) *MenuSectionService {
	baseService := NewGenericService[*domain.MenuSection](repo)
	return &MenuSectionService{genericService: baseService.(*genericService[*domain.MenuSection])}
}
