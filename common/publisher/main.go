package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/egigiffari/broker-go-rabbitmq/common"
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

	publisher := preparePublisher(rabbitMQ)

	go rabbitMQ.Reconnect(ctx, cancel)

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

			publisher.Publish(ctx, common.KEY, amqp.Publishing{
				Body:         []byte(message),
				DeliveryMode: amqp.Persistent,
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

func preparePublisher(rabbitMQ *pkg.RabbitMQ) *pkg.Publisher {
	publisher := pkg.Publisher{
		EventName:     common.EVENT,
		ConfirmNoWait: false,
		MQ:            rabbitMQ,
	}

	return &publisher
}
