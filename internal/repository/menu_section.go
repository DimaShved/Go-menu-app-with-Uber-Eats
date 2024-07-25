package repository

import (
	"gorm.io/gorm"
	"uber-go-menu-copy/internal/domain"
)

type MenuSectionRepo interface {
	IGenericRepository[*domain.MenuSections]
}

type menuSectionRepo struct {
	IGenericRepository[*domain.MenuSections]
}

func NewMenuSectionRepo(db *gorm.DB) MenuSectionRepo {
	return &menuSectionRepo{NewGenericRepo[*domain.MenuSections](db)}
}
