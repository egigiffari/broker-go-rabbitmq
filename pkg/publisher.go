package pkg

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	EventName     string
	ConfirmNoWait bool
	MQ            *RabbitMQ
	channel       *amqp.Channel
}

func (p *Publisher) Publish(ctx context.Context, key string, msg amqp.Publishing) error {
	channel, err := p.openChannel()
	if err != nil {
		return err
	}

	return channel.PublishWithContext(ctx, p.EventName, key, false, false, msg)
}

func (p *Publisher) PublishWithDeferedConfirm(ctx context.Context, key string, msg amqp.Publishing) (*amqp.DeferredConfirmation, error) {
	channel, err := p.openChannel()
	if err != nil {
		return nil, err
	}

	return channel.PublishWithDeferredConfirmWithContext(ctx, p.EventName, key, false, false, msg)
}

func (p *Publisher) openChannel() (*amqp.Channel, error) {
	if p.channel != nil && !p.channel.IsClosed() {
		return p.channel, nil
	}

	conn, err := p.MQ.Connection()
	if err != nil {
		return nil, err
	}

	log.Println("open new connection")
	newChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := newChannel.Confirm(p.ConfirmNoWait); err != nil {
		return nil, err
	}

	p.channel = newChannel
	return p.channel, nil
}
