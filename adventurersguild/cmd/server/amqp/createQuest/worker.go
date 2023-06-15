package createQuest

import (
	"context"
	"encoding/json"
	"log"

	"github.com/bocdagla/tavern/adventurersguild/cmd/server/amqp"
	"github.com/bocdagla/tavern/adventurersguild/cqrs/quest"
	"github.com/bocdagla/tavern/adventurersguild/cqrs/quest/command"
	"github.com/rabbitmq/amqp091-go"
)

// creates a worker that processes createQuest messages
type Worker struct {
	logger  *log.Logger
	handler quest.Handler
}

// creates a new Worker
func New(logger *log.Logger, handler quest.Handler) amqp.Worker {
	return &Worker{
		logger:  logger,
		handler: handler,
	}
}

// goroutine used to process the requests
func (c *Worker) Start(ctx context.Context, messages <-chan amqp091.Delivery) {
	var create *command.Create
	for delivery := range messages {
		c.logger.Printf("createQuest deliveryTag% v", delivery.DeliveryTag)
		err := json.Unmarshal(delivery.Body, create)
		if err != nil {
			c.logger.Printf("Failed to unmarshal request: %v", err)
		}
		err = c.handler.HandleCreate(ctx, *create)
		if err != nil {
			if err := delivery.Reject(false); err != nil {
				c.logger.Printf("Err delivery.Reject: %v", err)
			}
			c.logger.Printf("Failed to process createQuest: %v", err)
		} else {
			err = delivery.Ack(false)
			if err != nil {
				c.logger.Printf("Failed to acknowledge createQuest: %v", err)
			}
		}
	}

	c.logger.Printf("Deliveries channel closed")
}
