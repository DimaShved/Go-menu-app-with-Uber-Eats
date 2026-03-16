package menu_category

import "github.com/google/uuid"

type CreateRequest struct {
	SectionID   uuid.UUID `json:"section_id" validate:"required"`
	Name        string    `json:"name" validate:"required,min=2,max=255"`
	Description string    `json:"description"`
	IsAvailable bool      `json:"is_available" validate:"boolean,omitempty"`
}

type UpdateRequest = CreateRequest
