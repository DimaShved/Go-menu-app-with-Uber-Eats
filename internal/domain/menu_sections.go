package domain

import (
	"time"

	"github.com/google/uuid"
)

type MenuSection struct {
	ID                 uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	RestaurantID       uuid.UUID         `gorm:"type:uuid;not null" json:"restaurant_id" validate:"required,uuid"`
	Name               string            `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=2,max=255"`
	IsAvailable        bool              `gorm:"type:boolean;default:false" json:"is_available" validate:"boolean,omitempty"`
	CreatedAt          time.Time         `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt          time.Time         `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt          *time.Time        `gorm:"index" json:"deleted_at,omitempty"`
	Restaurant         Restaurant        `gorm:"foreignKey:RestaurantID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"restaurant" validate:"-"`
	MenuAvailabilityID *uuid.UUID        `gorm:"type:uuid;index;" json:"menu_availability_id,omitempty"`
	MenuAvailability   *MenuAvailability `gorm:"foreignKey:MenuAvailabilityID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"menu_availability,omitempty"`
}

func (ms MenuSection) GetID() uuid.UUID {
	return ms.ID
}

func (ms MenuSection) GetDeletedAt() *time.Time {
	return ms.DeletedAt
}
