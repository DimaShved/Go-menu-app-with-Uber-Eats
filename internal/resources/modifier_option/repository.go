package modifier_option

import (
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.ModifierOption]
}

func NewRepository() *Repository {
	return &Repository{GormRepository: crud.NewGormRepository[domain.ModifierOption]()}
}
