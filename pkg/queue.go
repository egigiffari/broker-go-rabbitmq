package pkg

import amqp "github.com/rabbitmq/amqp091-go"

type Queue struct {
	Name string
	Args amqp.Table
}

func (q *Queue) Declare(channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(q.Name, true, false, false, false, q.Args)
	return err
}

func (q *Queue) BindOnExchange(channel *amqp.Channel, key string, exchange string, args amqp.Table) error {
	return channel.QueueBind(q.Name, key, exchange, false, args)
}
