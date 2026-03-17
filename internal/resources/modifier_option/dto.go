package modifier_option

import "github.com/google/uuid"

type CreateRequest struct {
	Name         string    `json:"name" validate:"required"`
	Price        int       `json:"price" validate:"required"`
	MaxSelection int       `json:"max_selection" validate:"required"`
	IsAvailable  bool      `json:"is_available" validate:"boolean,omitempty"`
	ModifierID   uuid.UUID `json:"modifier_id" validate:"required"`
}

type UpdateRequest = CreateRequest
