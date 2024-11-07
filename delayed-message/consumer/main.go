package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	common "github.com/egigiffari/broker-go-rabbitmq/delayed-message"
	"github.com/egigiffari/broker-go-rabbitmq/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.TODO(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	rabbitMQ := getMQ()

	if err := common.Boot(rabbitMQ); err != nil {
		panic(err)
	}

	fmt.Println("====================")
	fmt.Println("Consume your message")
	fmt.Println("====================")

	consumer := getConsumer(rabbitMQ)
	var wg sync.WaitGroup
	wg.Add(1)

	go rabbitMQ.Reconnect(ctx, cancel)
	go consumer.Listen(ctx, &wg)

	<-ctx.Done()
	cancel()
	rabbitMQ.Close()
	wg.Wait()
	fmt.Println("goodbye")
}

func getMQ() *pkg.RabbitMQ {
	config := pkg.Config{
		Host:              "localhost",
		Port:              5672,
		Username:          "guest",
		Password:          "guest",
		MaxReconnect:      5,
		ReconnectInterval: 5,
	}

	return pkg.NewRabbitMQ(&config)
}

func prepareQueue(channel *amqp.Channel) error {
	queue := pkg.Queue{
		Name:    common.QUEUE,
		Durable: true,
		Args:    amqp.Table{},
	}

	if err := queue.Declare(channel); err != nil {
		return err
	}

	return queue.BindOnExchange(channel, "", common.EVENT, amqp.Table{})
}

func getConsumer(rabbitmq *pkg.RabbitMQ) *pkg.Consumer {
	return &pkg.Consumer{
		Name:      common.CONSUMER,
		Queue:     common.QUEUE,
		AuthAck:   true,
		NumWorker: 1,
		RabbitMQ:  rabbitmq,
		Handler: func(message amqp.Delivery) {
			fmt.Printf("recieved message: %s\n", string(message.Body))
		},
	}
}
