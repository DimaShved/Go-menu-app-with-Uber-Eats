package menu_section

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	RestaurantID       uuid.UUID  `json:"restaurant_id" validate:"required"`
	Name               string     `json:"name" validate:"required,min=2,max=255"`
	IsAvailable        bool       `json:"is_available" validate:"boolean,omitempty"`
	MenuAvailabilityID *uuid.UUID `json:"menu_availability_id,omitempty"`
}

type UpdateRequest = CreateRequest

type Response struct {
	ID                 uuid.UUID                 `json:"id,omitempty"`
	RestaurantID       uuid.UUID                 `json:"restaurant_id"`
	Name               string                    `json:"name"`
	IsAvailable        bool                      `json:"is_available"`
	CreatedAt          time.Time                 `json:"created_at,omitempty"`
	UpdatedAt          time.Time                 `json:"updated_at,omitempty"`
	Restaurant         RestaurantResponse        `json:"restaurant"`
	MenuAvailabilityID *uuid.UUID                `json:"menu_availability_id,omitempty"`
	MenuAvailability   *MenuAvailabilityResponse `json:"menu_availability,omitempty"`
}

type RestaurantResponse struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type MenuAvailabilityResponse struct {
	ID            uuid.UUID `json:"id,omitempty"`
	MenuSectionID uuid.UUID `json:"menu_section_id"`
	DayOfWeek     int       `json:"day_of_week"`
	OpenTime      int       `json:"open_time"`
	CloseTime     int       `json:"close_time"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
