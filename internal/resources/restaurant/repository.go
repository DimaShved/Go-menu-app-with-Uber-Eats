package restaurant

import (
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.Restaurant]
}

func NewRepository() *Repository {
	return &Repository{GormRepository: crud.NewGormRepository[domain.Restaurant]()}
}
