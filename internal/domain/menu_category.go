package domain

import (
	"github.com/google/uuid"
	"time"
)

type MenuCategory struct {
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id,omitempty"`
	SectionID   uuid.UUID    `gorm:"type:uuid;not null" json:"section_id" validate:"required,uuid"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=2,max=255"`
	Description string       `gorm:"type:varchar(255)" json:"description"`
	IsAvailable bool         `gorm:"type:boolean;default:true" json:"is_available" validate:"boolean"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   *time.Time   `gorm:"index" json:"deleted_at,omitempty"`
	Section     MenuSections `gorm:"foreignKey:SectionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"section" validate:"-"`
}

func (mc *MenuCategory) PreloadRelations() []string {
	return []string{"Section", "Section.Restaurant"}
}

func (mc *MenuCategory) GetID() uuid.UUID {
	return mc.ID
}
