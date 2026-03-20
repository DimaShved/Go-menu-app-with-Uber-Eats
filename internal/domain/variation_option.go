package domain

import (
	"github.com/google/uuid"
	"time"
)

type VariationOption struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Name        string     `gorm:"type:varchar(255);not null" json:"name"`
	Price       int        `gorm:"type:int;not null" json:"price"`
	IsAvailable bool       `gorm:"type:boolean;default:false" json:"is_available" validate:"boolean,omitempty"`
	VariationID uuid.UUID  `gorm:"type:uuid;not null" json:"variation_id"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

func (vo *VariationOption) GetID() uuid.UUID {
	return vo.ID
}
