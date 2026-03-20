package domain

import (
	"github.com/google/uuid"
	"time"
)

type Variation struct {
	ID         uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Name       string            `gorm:"type:varchar(255);not null" json:"name"`
	Options    []VariationOption `gorm:"foreignKey:VariationID" json:"options"`
	CategoryID uuid.UUID         `gorm:"type:uuid;not null" json:"category_id"`
	Category   MenuCategory      `gorm:"foreignKey:CategoryID" json:"-"`
	CreatedAt  time.Time         `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt  time.Time         `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt  *time.Time        `gorm:"index" json:"deleted_at,omitempty"`
}

func (v *Variation) GetID() uuid.UUID {
	return v.ID
}
