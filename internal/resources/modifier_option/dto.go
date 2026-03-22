package modifier_option

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	Name         string    `json:"name" validate:"required"`
	Price        int       `json:"price" validate:"required"`
	MaxSelection int       `json:"max_selection" validate:"required"`
	IsAvailable  bool      `json:"is_available" validate:"boolean,omitempty"`
	ModifierID   uuid.UUID `json:"modifier_id" validate:"required"`
}

type UpdateRequest = CreateRequest

type Response struct {
	ID           uuid.UUID `json:"id,omitempty"`
	Name         string    `json:"name"`
	Price        int       `json:"price"`
	MaxSelection int       `json:"max_selection"`
	IsAvailable  bool      `json:"is_available"`
	ModifierID   uuid.UUID `json:"modifier_id"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}
