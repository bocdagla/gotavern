// handler used to process the commands for the quests entity
package quest

import (
	"context"
	"log"

	"github.com/bocdagla/tavern/adventurersguild/cqrs/quest/command"
	"github.com/bocdagla/tavern/adventurersguild/pkg/domain/quest"
)

type Handler interface {
	HandleCreate(ctx context.Context, c command.Create) error
	HandleSetValid(ctx context.Context, v command.SetValid) error
}

type handlerImpl struct {
	logger *log.Logger
	rp     quest.Repository
}

func NewHandler(logger *log.Logger, rp quest.Repository) Handler {
	return &handlerImpl{
		logger: logger,
		rp:     rp,
	}
}

func (h *handlerImpl) HandleCreate(ctx context.Context, c command.Create) error {
	q, err := quest.NewQuest(c.Holder, c.Description)
	if err != nil {
		h.logger.Printf("An error ocurred when creating a quest: %v", err)
		return err
	}
	h.rp.Add(q)
	return nil
}

func (h *handlerImpl) HandleSetValid(ctx context.Context, v command.SetValid) error {
	return nil
}
