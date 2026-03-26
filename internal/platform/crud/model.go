package crud

import (
	"time"

	"github.com/google/uuid"
)

const (
	idColumn        = "id"
	deletedAtColumn = "deleted_at"
)

type SoftDeleteEntity interface {
	GetID() uuid.UUID
	GetDeletedAt() *time.Time
}
