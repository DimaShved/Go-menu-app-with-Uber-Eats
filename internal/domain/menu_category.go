package domain

import (
	"github.com/google/uuid"
	"time"
)

type MenuCategory struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	SectionID   uuid.UUID   `gorm:"type:uuid;not null" json:"section_id" validate:"required,uuid"`
	Name        string      `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=2,max=255"`
	Description string      `gorm:"type:varchar(255)" json:"description"`
	IsAvailable bool        `gorm:"type:boolean;default:false" json:"is_available" validate:"boolean,omitempty"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   *time.Time  `gorm:"index" json:"deleted_at,omitempty"`
	Section     MenuSection `gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"section" validate:"-"`
	Items       []MenuItem  `gorm:"many2many:item_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items" validate:"-"`
}

func (mc *MenuCategory) GetID() uuid.UUID {
	return mc.ID
}
