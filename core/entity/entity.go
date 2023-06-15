package entity

import "github.com/google/uuid"

type Entity interface {
	GetId() uuid.UUID
}

type EntityImpl struct {
	ID uuid.UUID
}

func (e *EntityImpl) GetId() uuid.UUID {
	return e.ID
}
