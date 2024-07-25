package domain

import (
	"github.com/google/uuid"
	"time"
)

type MenuSections struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id,omitempty"`
	RestaurantID uuid.UUID  `gorm:"type:uuid;not null" json:"restaurant_id" validate:"required,uuid"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=2,max=255"`
	IsAvailable  bool       `gorm:"type:boolean;default:true" json:"is_available" validate:"boolean"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
	Restaurant   Restaurant `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"restaurant" validate:"-"`
}

func (ms *MenuSections) PreloadRelations() []string {
	return []string{"Restaurant"}
}

func (ms *MenuSections) GetID() uuid.UUID {
	return ms.ID
}
