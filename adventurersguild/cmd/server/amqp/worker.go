package amqp

import (
	"context"

	"github.com/rabbitmq/amqp091-go"
)

// goroutine used to process the requests
type Worker interface {
	Start(ctx context.Context, messages <-chan amqp091.Delivery)
}
