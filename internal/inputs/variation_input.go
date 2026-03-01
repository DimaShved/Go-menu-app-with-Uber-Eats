package inputs

import "github.com/google/uuid"

type VariationInput struct {
	Name       string    `json:"name" validate:"required,min=2,max=255"`
	CategoryID uuid.UUID `json:"category_id" validate:"required"`
	Options    []struct {
		Name        string `json:"name" validate:"required"`
		Price       int    `json:"price" validate:"required"`
		IsAvailable bool   `json:"is_available" validate:"omitempty"`
	} `json:"options" validate:"omitempty,min=1,dive"`
}
