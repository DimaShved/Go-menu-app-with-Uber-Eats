package repository

import (
	"gorm.io/gorm"
	"uber-go-menu-copy/internal/domain"
)

type MenuCategoryRepo interface {
	IGenericRepository[*domain.MenuCategory]
}

type menuCategoryRepo struct {
	IGenericRepository[*domain.MenuCategory]
}

func NewMenuCategoryRepo(db *gorm.DB) MenuCategoryRepo {
	return &menuCategoryRepo{NewGenericRepo[*domain.MenuCategory](db)}
}
