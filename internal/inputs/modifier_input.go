package inputs

import "github.com/google/uuid"

type ModifierInput struct {
	Name              string    `json:"name" validate:"required,min=2,max=255"`
	TotalMaxSelection int       `json:"total_max_selection" validate:"required"`
	CategoryID        uuid.UUID `json:"category_id" validate:"required"`
	Options           []struct {
		Name         string `json:"name" validate:"required"`
		Price        int    `json:"price" validate:"required"`
		MaxSelection int    `json:"max_selection" validate:"required"`
		IsAvailable  bool   `json:"is_available" validate:"omitempty"`
	} `json:"options" validate:"omitempty,min=1,dive"`
}
