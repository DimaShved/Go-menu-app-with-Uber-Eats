package menu_item

type CreateRequest struct {
	Name        string   `json:"name" validate:"required,min=2,max=255"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	IsAvailable bool     `json:"is_available" validate:"boolean,omitempty"`
	Categories  []string `json:"categories" validate:"omitempty,dive,uuid"`
}

type UpdateRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	IsAvailable bool   `json:"is_available" validate:"boolean,omitempty"`
}
