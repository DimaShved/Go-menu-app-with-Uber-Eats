package domain

import (
	"github.com/google/uuid"
	"time"
)

type MenuAvailability struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	MenuSectionId uuid.UUID  `gorm:"type:uuid;not null;index:idx_menu_section,unique" json:"menu_section_id" validate:"required"`
	DayOfWeek     int        `gorm:"type:smallint;not null;index:idx_menu_section,unique" json:"day_of_week" validate:"required,oneof=1 2 3 4 5 6 7"`
	OpenTime      int        `gorm:"type:int;not null" json:"open_time" validate:"gte=0,lte=1439"`
	CloseTime     int        `gorm:"type:int;not null" json:"close_time" validate:"gte=0,lte=1439"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt     *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (ma *MenuAvailability) PreloadRelations() []string {
	return []string{}
}

func (ma *MenuAvailability) GetID() uuid.UUID {
	return ma.ID
}
