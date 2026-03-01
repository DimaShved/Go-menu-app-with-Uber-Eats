package inputs

import "github.com/google/uuid"

type MenuAvailabilityInput struct {
	MenuSectionId  uuid.UUID `json:"menu_section_id" validate:"required"`
	Availabilities []struct {
		DayOfWeek int `json:"day_of_week" validate:"required,oneof=1 2 3 4 5 6 7"`
		OpenTime  int `json:"open_time" validate:"gte=0,lte=1439"`
		CloseTime int `json:"close_time" validate:"gte=0,lte=1439"`
	} `json:"availabilities" validate:"required,dive"`
}
