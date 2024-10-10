package pkg

import amqp "github.com/rabbitmq/amqp091-go"

type Exchange struct {
	Name string
	Kind string
	Args amqp.Table
}

func (e *Exchange) Declare(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(e.Name, e.Kind, false, false, false, false, e.Args)
}
