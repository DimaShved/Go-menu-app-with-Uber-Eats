package menu_section

import (
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.MenuSection]
}

func NewRepository() *Repository {
	return &Repository{
		GormRepository: crud.NewGormRepository[domain.MenuSection](crud.QueryOptions{
			Preloads: []string{"Restaurant"},
		}),
	}
}
