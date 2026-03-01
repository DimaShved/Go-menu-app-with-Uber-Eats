package repository

import (
	"gorm.io/gorm"
	"uber-go-menu/internal/domain"
)

type MenuSectionRepo interface {
	IGenericRepository[*domain.MenuSection]
}

type menuSectionRepo struct {
	IGenericRepository[*domain.MenuSection]
}

func NewMenuSectionRepo(db *gorm.DB) MenuSectionRepo {
	return &menuSectionRepo{NewGenericRepo[*domain.MenuSection](db)}
}
