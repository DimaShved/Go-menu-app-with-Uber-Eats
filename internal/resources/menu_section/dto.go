package menu_section

import "github.com/google/uuid"

type CreateRequest struct {
	RestaurantID       uuid.UUID  `json:"restaurant_id" validate:"required"`
	Name               string     `json:"name" validate:"required,min=2,max=255"`
	IsAvailable        bool       `json:"is_available" validate:"boolean,omitempty"`
	MenuAvailabilityID *uuid.UUID `json:"menu_availability_id,omitempty"`
}

type UpdateRequest = CreateRequest
