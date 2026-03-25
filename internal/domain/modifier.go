package domain

import (
	"github.com/google/uuid"
	"time"
)

type Modifier struct {
	ID                uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id,omitempty"`
	Name              string           `gorm:"type:varchar(255);not null" json:"name"`
	TotalMaxSelection int              `gorm:"not null" json:"total_max_selection"`
	Options           []ModifierOption `gorm:"foreignKey:ModifierID" json:"options"`
	CategoryID        uuid.UUID        `gorm:"type:uuid;not null" json:"category_id"`
	CreatedAt         time.Time        `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt         time.Time        `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt         *time.Time       `gorm:"index" json:"deleted_at,omitempty"`
}

func (m Modifier) GetID() uuid.UUID {
	return m.ID
}

func (m Modifier) GetDeletedAt() *time.Time {
	return m.DeletedAt
}
