package menu_availability

import (
	"time"

	"github.com/google/uuid"
)

type AvailabilityRequest struct {
	DayOfWeek int `json:"day_of_week" validate:"required,oneof=1 2 3 4 5 6 7"`
	OpenTime  int `json:"open_time" validate:"gte=0,lte=1439"`
	CloseTime int `json:"close_time" validate:"gte=0,lte=1439"`
}

type CreateRequest struct {
	MenuSectionID uuid.UUID `json:"menu_section_id" validate:"required"`
	DayOfWeek     int       `json:"day_of_week" validate:"required,oneof=1 2 3 4 5 6 7"`
	OpenTime      int       `json:"open_time" validate:"gte=0,lte=1439"`
	CloseTime     int       `json:"close_time" validate:"gte=0,lte=1439"`
}

type BatchUpsertRequest struct {
	MenuSectionID  uuid.UUID             `json:"menu_section_id" validate:"required"`
	Availabilities []AvailabilityRequest `json:"availabilities" validate:"required,min=1,dive"`
}

type UpdateRequest = CreateRequest

type Response struct {
	ID            uuid.UUID `json:"id,omitempty"`
	MenuSectionID uuid.UUID `json:"menu_section_id"`
	DayOfWeek     int       `json:"day_of_week"`
	OpenTime      int       `json:"open_time"`
	CloseTime     int       `json:"close_time"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}
