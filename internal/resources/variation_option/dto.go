package variation_option

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	Name        string    `json:"name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	IsAvailable bool      `json:"is_available" validate:"boolean,omitempty"`
	VariationID uuid.UUID `json:"variation_id" validate:"required"`
}

type UpdateRequest = CreateRequest

type Response struct {
	ID          uuid.UUID `json:"id,omitempty"`
	Name        string    `json:"name"`
	Price       int       `json:"price"`
	IsAvailable bool      `json:"is_available"`
	VariationID uuid.UUID `json:"variation_id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}
