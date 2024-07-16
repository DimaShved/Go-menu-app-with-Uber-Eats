package service

import (
	"uber-go-menu-copy/internal/domain"
	"uber-go-menu-copy/internal/repository"
)

type RestaurantService struct {
	*genericService[*domain.Restaurant]
}

func NewRestaurantService(repo repository.IGenericRepository[*domain.Restaurant]) *RestaurantService {
	baseService := NewGenericService[*domain.Restaurant](repo)
	return &RestaurantService{genericService: baseService.(*genericService[*domain.Restaurant])}
}
