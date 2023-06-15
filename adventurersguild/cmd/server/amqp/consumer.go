package amqp

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/rabbitmq/amqp091-go"
)

const (
	exchangeKind       = "direct"
	exchangeDurable    = true
	exchangeAutoDelete = false
	exchangeInternal   = false
	exchangeNoWait     = false

	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	prefetchCount  = 1
	prefetchSize   = 0
	prefetchGlobal = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type Consumer interface {
	CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp091.Channel, error)
	StartConsumer(workerPoolSize int, queueName, consumerTag string, ch *amqp091.Channel) error
}

type questConsumer struct {
	amqpConn *amqp091.Connection
	logger   *log.Logger
	worker   Worker
}

func New(conn *amqp091.Connection, logger *log.Logger, worker Worker) Consumer {
	return &questConsumer{
		amqpConn: conn,
		logger:   logger,
		worker:   worker,
	}
}

// Creates a chanel
func (c *questConsumer) CreateChannel(exchangeName, queueName, bindingKey, consumerTag string) (*amqp091.Channel, error) {
	ch, err := c.amqpConn.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "Error amqpConn.Channel")
	}

	c.logger.Printf("Declaring exchange: %s", exchangeName)
	err = ch.ExchangeDeclare(
		exchangeName,
		exchangeKind,
		exchangeDurable,
		exchangeAutoDelete,
		exchangeInternal,
		exchangeNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.ExchangeDeclare")
	}

	queue, err := ch.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueDeclare")
	}

	c.logger.Printf("Declared queue, binding it to exchange: Queue: %v, messagesCount: %v, "+
		"consumerCount: %v, exchange: %v, bindingKey: %v",
		queue.Name,
		queue.Messages,
		queue.Consumers,
		exchangeName,
		bindingKey,
	)

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		queueNoWait,
		nil,
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error ch.QueueBind")
	}

	c.logger.Printf("Queue bound to exchange, starting to consume from queue, consumerTag: %v", consumerTag)

	err = ch.Qos(
		prefetchCount,  // prefetch count
		prefetchSize,   // prefetch size
		prefetchGlobal, // global
	)
	if err != nil {
		return nil, errors.Wrap(err, "Error  ch.Qos")
	}

	return ch, nil
}

// StartConsumer Start new rabbitmq consumer
func (c *questConsumer) StartConsumer(workerPoolSize int, queueName, consumerTag string, ch *amqp091.Channel) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer ch.Close()

	deliveries, err := ch.Consume(
		queueName,
		consumerTag,
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "Consume")
	}

	for i := 0; i < workerPoolSize; i++ {
		go c.worker.Start(ctx, deliveries)
	}

	chanErr := <-ch.NotifyClose(make(chan *amqp091.Error))
	c.logger.Printf("ch.NotifyClose: %v", chanErr)
	return chanErr
}
