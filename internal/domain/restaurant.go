package domain

import (
	"github.com/google/uuid"
	"time"
)

type Restaurant struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Name      string     `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=2,max=255"`
	Address   string     `gorm:"type:varchar(255);not null" json:"address" validate:"required,min=2,max=255"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (r Restaurant) GetID() uuid.UUID {
	return r.ID
}

func (r Restaurant) GetDeletedAt() *time.Time {
	return r.DeletedAt
}
