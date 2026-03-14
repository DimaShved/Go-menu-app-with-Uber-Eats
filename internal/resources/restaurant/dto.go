package restaurant

type CreateRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Address string `json:"address" validate:"required,min=2,max=255"`
}

type UpdateRequest = CreateRequest
