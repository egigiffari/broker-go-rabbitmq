package pkg

import (
	"context"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageHandler = func(message amqp.Delivery)

type Consumer struct {
	Name      string
	Queue     string
	AuthAck   bool
	NumWorker int
	RabbitMQ  *RabbitMQ
	Handler   MessageHandler
}

func (c *Consumer) Listen(ctx context.Context, wg *sync.WaitGroup) {
	for {
		for i := 0; i < c.NumWorker; i++ {
			fmt.Printf("consumer running worker %d\n", i+1)
			go c.Worker(ctx, c.Handler)
		}

		reconnect := <-c.RabbitMQ.NotifyReconnect(make(chan bool))
		fmt.Printf("reconnect: %v\n", reconnect)
		if reconnect {
			continue
		}

		break
	}

	<-ctx.Done()
	wg.Done()
	fmt.Println("success stop consumer")
}

func (c *Consumer) Worker(ctx context.Context, handler MessageHandler) {
	messages, err := c.Messages(ctx)
	if err != nil {
		fmt.Printf("failed connect consumer err: %s\n", err.Error())
		return
	}

	fmt.Println("consumer connected")
	for message := range messages {
		handler(message)
	}

	fmt.Println("stoped consumer")
}

func (c *Consumer) Messages(ctx context.Context) (<-chan amqp.Delivery, error) {
	connection, err := c.RabbitMQ.Connection()
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return channel.ConsumeWithContext(ctx, c.Queue, c.Name, c.AuthAck, false, false, false, amqp.Table{})
}
