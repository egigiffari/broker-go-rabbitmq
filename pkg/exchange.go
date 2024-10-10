package pkg

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Exchange struct {
	Name     string
	Kind     string
	Durable  bool
	Args     amqp.Table
	RabbitMQ *RabbitMQ
}

func (e *Exchange) Declare() error {
	channel, err := e.channel()
	if err != nil {
		return err
	}

	return channel.ExchangeDeclare(e.Name, e.Kind, e.Durable, false, false, false, e.Args)
}

func (e *Exchange) Publish(ctx context.Context, key string, mandatory bool, imediate bool, msg amqp.Publishing) error {
	channel, err := e.channel()
	if err != nil {
		return err
	}

	return channel.PublishWithContext(ctx, e.Name, key, mandatory, imediate, msg)
}

func (e *Exchange) channel() (*amqp.Channel, error) {
	connection, err := e.RabbitMQ.Connection()
	if err != nil {
		return nil, err
	}

	return connection.Channel()
}
