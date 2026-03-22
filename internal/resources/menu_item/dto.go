package menu_item

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	Name        string   `json:"name" validate:"required,min=2,max=255"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	IsAvailable bool     `json:"is_available" validate:"boolean,omitempty"`
	Categories  []string `json:"categories" validate:"omitempty,dive,uuid"`
}

type UpdateRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	IsAvailable bool   `json:"is_available" validate:"boolean,omitempty"`
}

type Response struct {
	ID          uuid.UUID          `json:"id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Price       int                `json:"price"`
	IsAvailable bool               `json:"is_available"`
	CreatedAt   time.Time          `json:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty"`
	Categories  []CategoryResponse `json:"categories"`
}

type CategoryResponse struct {
	ID          uuid.UUID       `json:"id,omitempty"`
	SectionID   uuid.UUID       `json:"section_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsAvailable bool            `json:"is_available"`
	CreatedAt   time.Time       `json:"created_at,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at,omitempty"`
	Section     SectionResponse `json:"section"`
}

type SectionResponse struct {
	ID                 uuid.UUID          `json:"id,omitempty"`
	RestaurantID       uuid.UUID          `json:"restaurant_id"`
	Name               string             `json:"name"`
	IsAvailable        bool               `json:"is_available"`
	CreatedAt          time.Time          `json:"created_at,omitempty"`
	UpdatedAt          time.Time          `json:"updated_at,omitempty"`
	Restaurant         RestaurantResponse `json:"restaurant"`
	MenuAvailabilityID *uuid.UUID         `json:"menu_availability_id,omitempty"`
}

type RestaurantResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
