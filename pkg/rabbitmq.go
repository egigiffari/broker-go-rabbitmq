package pkg

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	config           *Config
	connection       *amqp.Connection
	stop             bool
	reconnectChan    []chan bool
	attemptReconnect int
}

func NewRabbitMQ(config *Config) *RabbitMQ {
	return &RabbitMQ{
		config: config,
	}
}

func (r *RabbitMQ) Connection() (*amqp.Connection, error) {
	if r.connection == nil {
		return r.Connect()
	}

	return r.connection, nil
}

func (r *RabbitMQ) Connect() (*amqp.Connection, error) {
	conn, err := amqp.Dial(r.config.Dns())
	if err != nil {
		return nil, err
	}

	r.connection = conn
	r.attemptReconnect = 0
	fmt.Println("success connect")
	return r.connection, nil
}

func (r *RabbitMQ) Reconnect(ctx context.Context, stop context.CancelFunc) {
	for {
		select {
		case <-ctx.Done():
			return
		case notify := <-r.connection.NotifyClose(make(chan *amqp.Error)):
			fmt.Printf("close cause code: %d, reason %s\n", notify.Code, notify.Reason)

			for {
				if r.attemptReconnect > r.config.MaxReconnect {
					fmt.Println("rabbitmq has reached max reconnect attempt")
					stop()
					return
				}

				fmt.Printf("rabbitmq reconnect %d/%d\n", r.attemptReconnect, r.config.MaxReconnect)
				if _, err := r.Connect(); err == nil {
					r.sendReconnect(true)
					break
				}

				time.Sleep(time.Duration(r.config.ReconnectInterval) * time.Second)
				r.attemptReconnect += 1
			}
		}
	}
}

func (r *RabbitMQ) IsStop() bool {
	return r.stop
}

func (r *RabbitMQ) NotifyReconnect(reciever chan bool) chan bool {
	if r.stop {
		close(reciever)
	} else {
		r.reconnectChan = append(r.reconnectChan, reciever)
	}
	return reciever
}

func (r *RabbitMQ) Close() error {
	r.stop = true
	r.sendReconnect(false)
	return r.connection.Close()
}

func (r *RabbitMQ) sendReconnect(state bool) {
	for _, rc := range r.reconnectChan {
		rc <- state
	}

	for _, rc := range r.reconnectChan {
		close(rc)
	}

	r.reconnectChan = nil
}
