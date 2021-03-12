package checker

import (
	"github.com/yanpin-dev/propeller/pkg/health"
	"github.com/streadway/amqp"
)

type RabbitChecker struct {
	name string
	conn *amqp.Connection
}

func NewRabbitChecker(conn *amqp.Connection) health.Checker {
	return RabbitChecker{
		name: "rabbit",
		conn: conn,
	}
}

func (c RabbitChecker) Check() health.Health {
	result := health.NewHealth()
	channel, err := c.conn.Channel()
	if err != nil {
		result.Down().AddInfo("error", err.Error())
		return result
	}
	defer channel.Close()
	result.Up()
	return result
}

func (c RabbitChecker) Name() string {
	return c.name
}
