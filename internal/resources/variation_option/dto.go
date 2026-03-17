package variation_option

import "github.com/google/uuid"

type CreateRequest struct {
	Name        string    `json:"name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	IsAvailable bool      `json:"is_available" validate:"boolean,omitempty"`
	VariationID uuid.UUID `json:"variation_id" validate:"required"`
}

type UpdateRequest = CreateRequest
