package pkg

import "fmt"

type Config struct {
	Host              string
	Port              int
	Username          string
	Password          string
	MaxReconnect      int
	ReconnectInterval int
}

func (c *Config) Dns() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
}
