# golang-rabbimq
simulate how to golang handle publish or consume data from rabbitmq, i hope this project can help you


## How to run
this project require `docker` for running rabbitmq. running docker with `docker-compose.yml` for configuration, after container `golang-rabbitmq` is running you can run sample usecase.

### usecase

#### common
this the common usecase for rabbitqm as message queue, publisher for publish message and consumer for recieved message from publisher. for start usecase make sure you start publisher first before consumer

#### delayed-message
if publisher has somecase for delayed message before consumer recieved message depend on duration, you must install `rabbitmq_delayed_message_exchange` first, you can get more info on [rabbitmq doc](https://www.rabbitmq.com/blog/2015/04/16/scheduling-messages-with-rabbitmq). on this case publisher will push message after 5 second

```go
exchange := pkg.Exchange{
		Name:     EXCHANGE,
		Kind:     "x-delayed-message",
		Durable:  true,
		Args:     amqp.Table{"x-delayed-type": "direct"},
		RabbitMQ: rabbitMQ,
	}
```
> Note: `x-delayed-message` is required for declare exchange delayed-message.


```go
exchange.Publish(ctx, "", false, false, amqp.Publishing{
				Body:         []byte(message),
				DeliveryMode: amqp.Persistent,
				Headers: amqp.Table{
					"x-delay": 5000, // in millisecond
				},
			})
```
> Note: `x-delay` is required for set duration delay before publish message to queue
