// Package used to get data from the repository
package quest

import (
	"github.com/bocdagla/tavern/adventurersguild/cqrs/quest/query"
	"github.com/bocdagla/tavern/adventurersguild/pkg/domain/quest"
	"github.com/google/uuid"
)

type QuestDisplay struct {
	Id            uuid.UUID
	Rank          int8
	Certification uuid.UUID
	Description   string
	Holder        uuid.UUID
}

func fromQuestAgregate(q quest.Quest) QuestDisplay {
	return QuestDisplay{
		q.ID,
		int8(q.Rank),
		q.Certification,
		q.Description,
		q.Holder,
	}
}

type Dao interface {
	GetById(q query.ById) (QuestDisplay, error)
	GetAll(q query.All) ([]QuestDisplay, error)
}

type daoImpl struct {
	rp quest.Repository
}

func NewDao(rp quest.Repository) Dao {
	return &daoImpl{
		rp: rp,
	}
}

func (d *daoImpl) GetById(q query.ById) (QuestDisplay, error) {
	dquest, err := d.rp.Get(q.Id)
	if err != nil {
		return QuestDisplay{}, err
	}
	return fromQuestAgregate(dquest), nil
}

func (d *daoImpl) GetAll(q query.All) ([]QuestDisplay, error) {
	dquests, err := d.rp.GetAll()
	if err != nil {
		return []QuestDisplay{}, err
	}

	var result []QuestDisplay
	for _, dquest := range dquests {
		result = append(result, fromQuestAgregate(dquest))
	}
	return result, nil
}
