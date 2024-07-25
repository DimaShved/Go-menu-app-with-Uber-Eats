package domain

import (
	"github.com/google/uuid"
	"time"
)

type MenuItem struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id,omitempty"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name" validate:"required,min=2,max=255"`
	Description string         `gorm:"type:varchar(255)" json:"description"`
	Price       int            `gorm:"type:int" json:"price"`
	IsAvailable bool           `gorm:"type:boolean;default:true" json:"is_available" validate:"boolean"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt   *time.Time     `gorm:"index" json:"deleted_at,omitempty"`
	Categories  []MenuCategory `gorm:"many2many:item_categories;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"categories"`
}

func (mi *MenuItem) PreloadRelations() []string {
	return []string{"Categories", "Categories.Section", "Categories.Section.Restaurant"}
}

func (mi *MenuItem) GetID() uuid.UUID {
	return mi.ID
}
