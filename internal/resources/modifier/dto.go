package modifier

import (
	"time"

	"github.com/google/uuid"
)

type OptionRequest struct {
	Name         string `json:"name" validate:"required"`
	Price        int    `json:"price" validate:"required"`
	MaxSelection int    `json:"max_selection" validate:"required"`
	IsAvailable  bool   `json:"is_available" validate:"boolean,omitempty"`
}

type CreateRequest struct {
	Name              string          `json:"name" validate:"required,min=2,max=255"`
	TotalMaxSelection int             `json:"total_max_selection" validate:"required"`
	CategoryID        uuid.UUID       `json:"category_id" validate:"required"`
	Options           []OptionRequest `json:"options" validate:"omitempty,min=1,dive"`
}

type UpdateRequest struct {
	Name              string    `json:"name" validate:"required,min=2,max=255"`
	TotalMaxSelection int       `json:"total_max_selection" validate:"required"`
	CategoryID        uuid.UUID `json:"category_id" validate:"required"`
}

type Response struct {
	ID                uuid.UUID        `json:"id,omitempty"`
	Name              string           `json:"name"`
	TotalMaxSelection int              `json:"total_max_selection"`
	Options           []OptionResponse `json:"options"`
	CategoryID        uuid.UUID        `json:"category_id"`
	CreatedAt         time.Time        `json:"created_at,omitempty"`
	UpdatedAt         time.Time        `json:"updated_at,omitempty"`
}

type OptionResponse struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name"`
	Price        int       `json:"price"`
	MaxSelection int       `json:"max_selection"`
	IsAvailable  bool      `json:"is_available"`
	ModifierID   uuid.UUID `json:"modifier_id"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
