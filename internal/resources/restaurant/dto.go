package restaurant

import (
	"time"

	"github.com/google/uuid"
)

type CreateRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Address string `json:"address" validate:"required,min=2,max=255"`
}

type UpdateRequest = CreateRequest

type Response struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
