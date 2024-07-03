package interfaces

import "github.com/google/uuid"

type Identifiable interface {
	GetID() uuid.UUID
}
