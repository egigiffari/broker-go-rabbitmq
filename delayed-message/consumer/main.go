package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	"github.com/egigiffari/broker-go-rabbitmq/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EXCHANGE = "sample-delayed-exchange"
	QUEUE    = "sample-queue"
	CONSUMER = "sample-consumer"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.TODO(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	rabbitMQ := getMQ()

	connection, err := rabbitMQ.Connection()
	if err != nil {
		panic(fmt.Errorf("rabbitmq failed to connect err: %s", err.Error()))
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(fmt.Errorf("rabbitmq channel not response err: %s", err.Error()))
	}

	if err := prepareQueue(channel); err != nil {
		panic(fmt.Errorf("rabbitmq failed declare queue err: %s", err.Error()))
	}

	fmt.Println("====================")
	fmt.Println("Consume your message")
	fmt.Println("====================")

	consumer := getConsumer(rabbitMQ)
	var wg sync.WaitGroup
	wg.Add(1)

	go rabbitMQ.Reconnect(ctx)
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
		Name:    QUEUE,
		Durable: true,
		Args:    amqp.Table{},
	}

	if err := queue.Declare(channel); err != nil {
		return err
	}

	return queue.BindOnExchange(channel, "", EXCHANGE, amqp.Table{})
}

func getConsumer(rabbitmq *pkg.RabbitMQ) *pkg.Consumer {
	return &pkg.Consumer{
		Name:      CONSUMER,
		Queue:     QUEUE,
		AuthAck:   true,
		NumWorker: 1,
		RabbitMQ:  rabbitmq,
		Handler: func(message amqp.Delivery) {
			fmt.Printf("recieved message: %s\n", string(message.Body))
		},
	}
}
