// The quest package of the domain allows to perform operations with the Agregate root quest
package quest

import (
	"errors"

	"github.com/bocdagla/tavern/core/entity"
	"github.com/google/uuid"
)

// Rank of the quest, being S the highes grade and E tge lowes in this order S>A>B>C>D>E
type Rank int8

const (
	S Rank = iota
	A
	B
	C
	D
	E
)

var (
	//Error for when the quest holder is invalid
	ErrHolderInvalid = errors.New("the quest holder is invalid")
	//Error for when the description of the quest  is empty
	ErrDescriptionEmpty = errors.New("the description is empty")
)

// Entity defining a quest
type Quest struct {
	entity.EntityImpl
	Holder        uuid.UUID
	Description   string
	Rank          Rank
	Certification uuid.UUID
}

// Returns a new quest
// Throws an error if the description is empty
func NewQuest(holder uuid.UUID, description string) (Quest, error) {
	if description == "" {
		return Quest{}, ErrDescriptionEmpty
	}
	if holder == uuid.Nil {
		return Quest{}, ErrHolderInvalid
	}

	return Quest{
		Holder:      holder,
		Description: description,
	}, nil
}
