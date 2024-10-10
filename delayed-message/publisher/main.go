package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/egigiffari/broker-go-rabbitmq/pkg"
	amqp "github.com/rabbitmq/amqp091-go"
)

const EXCHANGE = "sample-delayed-exchange"

func main() {
	ctx, cancel := signal.NotifyContext(context.TODO(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	rabbitMQ := getMQ()
	connection, err := rabbitMQ.Connection()
	if err != nil {
		panic(fmt.Errorf("rabbitmq failed to connect err: %s", err.Error()))
	}
	defer connection.Close()

	exchange, err := prepareExchange(rabbitMQ)
	if err != nil {
		panic(fmt.Errorf("rabbitmq failed declare exchange err: %s", err.Error()))
	}

	go rabbitMQ.Reconnect(ctx)
	var message string
	fmt.Println("====================")
	fmt.Println("Publish your message")
	fmt.Println("====================")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("goodbye")
			return
		default:
			fmt.Print("pesan :")
			fmt.Scanf("%s", &message)
			if message == "" {
				continue
			}

			exchange.Publish(ctx, "", false, false, amqp.Publishing{
				Body:         []byte(message),
				DeliveryMode: amqp.Persistent,
				Headers: amqp.Table{
					"x-delay": 5000, // in millisecond
				},
			})
			message = ""
		}
	}
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

func prepareExchange(rabbitMQ *pkg.RabbitMQ) (*pkg.Exchange, error) {
	exchange := pkg.Exchange{
		Name:     EXCHANGE,
		Kind:     "x-delayed-message",
		Durable:  true,
		Args:     amqp.Table{"x-delayed-type": "direct"},
		RabbitMQ: rabbitMQ,
	}

	if err := exchange.Declare(); err != nil {
		return nil, err
	}

	return &exchange, nil
}
