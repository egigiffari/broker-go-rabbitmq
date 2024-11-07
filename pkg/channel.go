package pkg

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Channel struct {
	MQ *RabbitMQ
	ch *amqp.Channel
}

func NewChannel(mq *RabbitMQ) *Channel {
	return &Channel{MQ: mq}
}

func (c *Channel) ExchangeDeclare(name, kind string, durable, autoDelete bool, args amqp.Table) error {
	ch, err := c.open()
	if err != nil {
		return err
	}

	return ch.ExchangeDeclare(name, kind, durable, autoDelete, false, false, args)
}

func (c *Channel) QueueDeclare(name string, durable, autoDelete bool, args amqp.Table) error {
	ch, err := c.open()
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(name, durable, autoDelete, false, false, args)
	return err
}

func (c *Channel) QueueBind(name, key, exchange string, args amqp.Table) error {
	ch, err := c.open()
	if err != nil {
		return err
	}

	return ch.QueueBind(name, key, exchange, false, args)
}

func (c *Channel) Close() error {
	if c.ch == nil || c.ch.IsClosed() {
		return nil
	}

	return c.ch.Close()
}

func (c *Channel) open() (*amqp.Channel, error) {
	if c.ch != nil && !c.ch.IsClosed() {
		return c.ch, nil
	}

	conn, err := c.MQ.Connection()
	if err != nil {
		return nil, err
	}

	newChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	c.ch = newChannel
	return c.ch, nil
}
