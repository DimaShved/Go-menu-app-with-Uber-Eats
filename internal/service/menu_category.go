package service

import (
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/repository"
)

type MenuCategoryService struct {
	*genericService[*domain.MenuCategory]
}

func NewMenuCategoryService(repo repository.IGenericRepository[*domain.MenuCategory]) *MenuCategoryService {
	baseService := NewGenericService[*domain.MenuCategory](repo)
	return &MenuCategoryService{genericService: baseService.(*genericService[*domain.MenuCategory])}
}
