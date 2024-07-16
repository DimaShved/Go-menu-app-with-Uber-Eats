package repository

import (
	"gorm.io/gorm"
	"uber-go-menu-copy/internal/domain"
)

type RestaurantRepo interface {
	IGenericRepository[*domain.Restaurant]
}

type restaurantRepo struct {
	IGenericRepository[*domain.Restaurant]
}

func NewRestaurantRepo(db *gorm.DB) RestaurantRepo {
	return &restaurantRepo{NewGenericRepo[*domain.Restaurant](db)}
}
