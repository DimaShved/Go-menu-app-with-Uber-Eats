package modifier

import "github.com/google/uuid"

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
