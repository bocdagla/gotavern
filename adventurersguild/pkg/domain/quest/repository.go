package quest

import (
	"github.com/google/uuid"
)

type Repository interface {
	Get(id uuid.UUID) (Quest, error)
	GetAll() ([]Quest, error)
	Add(quest Quest) (Quest, error)
	Edit(quest Quest) (Quest, error)
	Delete(id uuid.UUID) error
}
