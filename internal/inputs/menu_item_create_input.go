package inputs

type MenuItemCreateInput struct {
	Name        string   `json:"name" validate:"required,min=2,max=255"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	IsAvailable bool     `json:"is_available"`
	Categories  []string `json:"categories"`
}
