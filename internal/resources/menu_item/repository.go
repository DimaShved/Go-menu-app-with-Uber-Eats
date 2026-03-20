package menu_item

import (
	"uber-go-menu/internal/domain"
	"uber-go-menu/internal/platform/crud"
)

type Repository struct {
	*crud.GormRepository[domain.MenuItem]
}

func NewRepository() *Repository {
	return &Repository{
		GormRepository: crud.NewGormRepository[domain.MenuItem](crud.QueryOptions{
			Preloads: []string{"Categories", "Categories.Section", "Categories.Section.Restaurant"},
		}),
	}
}
