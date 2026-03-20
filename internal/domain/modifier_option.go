package domain

import (
	"github.com/google/uuid"
	"time"
)

type ModifierOption struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Name         string     `gorm:"type:varchar(255);not null" json:"name"`
	Price        int        `gorm:"type:int;not null" json:"price"`
	MaxSelection int        `gorm:"not null" json:"max_selection"`
	IsAvailable  bool       `gorm:"type:boolean;default:false" json:"is_available" validate:"boolean,omitempty"`
	ModifierID   uuid.UUID  `gorm:"type:uuid;not null" json:"modifier_id"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt    *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (mo *ModifierOption) GetID() uuid.UUID {
	return mo.ID
}
