package common

import (
	"fmt"

	"github.com/egigiffari/broker-go-rabbitmq/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EVENT    = "sample-delayed-exchange"
	QUEUE    = "sample-queue"
	KEY      = ""
	CONSUMER = "sample-consumer"
)

func Boot(mq *pkg.RabbitMQ) error {
	channel := pkg.NewChannel(mq)

	err := channel.ExchangeDeclare(EVENT, "x-delayed-message", true, false, amqp.Table{"x-delayed-type": "direct"})
	if err != nil {
		return fmt.Errorf("failed to declare exchange, err: %s", err.Error())
	}

	err = channel.QueueDeclare(QUEUE, true, false, amqp.Table{})
	if err != nil {
		return fmt.Errorf("failed to declare queue, err: %s", err.Error())
	}

	err = channel.QueueBind(QUEUE, KEY, EVENT, amqp.Table{})
	if err != nil {
		return fmt.Errorf("failed to bind queue and exchange, err: %s", err.Error())
	}

	err = channel.Close()
	if err != nil {
		return fmt.Errorf("failed to close channel, err: %s", err.Error())
	}

	return nil
}
