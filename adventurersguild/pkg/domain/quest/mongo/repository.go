// This package contains operations to work with mongoDb
package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/bocdagla/tavern/adventurersguild/pkg/domain/quest"
	"github.com/bocdagla/tavern/core/db"
	"github.com/bocdagla/tavern/core/entity"
	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
)

// Mongo db repository for the Quest agregate
type Repository struct {
	quests db.Collection
}

// Creates a new Database conection and asigns it to the repository
func New(mc db.Collection) *Repository {
	return &Repository{
		quests: mc,
	}
}

// gets a quest by the id of the quest in the database
func (r *Repository) Get(id uuid.UUID) (quest.Quest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := r.quests.FindOne(ctx, bson.M{"_id": id})

	var c mongoQuest
	err := result.Decode(&c)
	if err != nil {
		return quest.Quest{}, err
	}
	return c.toAggregate(), nil
}

// gets all quests
func (r *Repository) GetAll() ([]quest.Quest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.quests.Find(ctx, nil)
	if err != nil {
		return []quest.Quest{}, err
	}

	var c []mongoQuest
	err = result.Decode(&c)
	if err != nil {
		return []quest.Quest{}, err
	}
	quests := make([]quest.Quest, 0)
	for _, q := range c {
		quests = append(quests, q.toAggregate())
	}
	return quests, nil
}

// Adds a new quest document to the quest database
func (r *Repository) Add(q quest.Quest) (quest.Quest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rq := newMongoQuestFromAggregate(q)
	result, err := r.quests.InsertOne(ctx, rq)
	if err != nil {
		return quest.Quest{}, err
	}
	if t, err := result.InsertedID.(uuid.UUID); err {
		return quest.Quest{}, errors.New("could not parse saved id as uuid.UUID")
	} else {
		return r.Get(t)
	}
}

// Edits an already existing quest document overriding it with the values passed
func (r *Repository) Edit(q quest.Quest) (quest.Quest, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rq := newMongoQuestFromAggregate(q)
	result, err := r.quests.UpdateByID(ctx, q.ID, rq)
	if err != nil {
		return quest.Quest{}, err
	}
	if t, err := result.UpsertedID.(uuid.UUID); err {
		return quest.Quest{}, errors.New("could not parse saved id as uuid.UUID")
	} else {
		return r.Get(t)
	}
}

// Deletes an exisiting quest with the Id passed by parameter
func (r *Repository) Delete(id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.quests.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

// Internal class to handle the conversion of the database into the domain model
type mongoQuest struct {
	ID            uuid.UUID `bson:"_id"`
	Description   string    `bson:"description"`
	Holder        uuid.UUID `bson:"holder"`
	Rank          int8      `bson:"rank"`
	Certification uuid.UUID `bson:"certification"`
}

// mapper to convert Entity object to database object
func newMongoQuestFromAggregate(q quest.Quest) mongoQuest {
	return mongoQuest{
		ID:            q.ID,
		Description:   q.Description,
		Holder:        q.Holder,
		Rank:          int8(q.Rank),
		Certification: q.Certification,
	}
}

// mapper to convert database object to database entity object
func (q *mongoQuest) toAggregate() quest.Quest {
	return quest.Quest{
		EntityImpl: entity.EntityImpl{
			ID: q.ID,
		},
		Holder:        q.Holder,
		Description:   q.Description,
		Rank:          quest.Rank(q.Rank),
		Certification: q.Certification,
	}
}
