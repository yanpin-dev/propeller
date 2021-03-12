package rabbit

import (
	"fmt"
	"github.com/streadway/amqp"
	"net/url"
)

const (
	DelayedMessage   = "x-delayed-message"
	DelayedType      = "x-delayed-type"
	DefaultDelayType = "topic"
)

type Client interface {
	CreateScheme(s *Options) error
	SendMessage(ex, key string, data []byte, headers map[string]interface{}) error
	ProcessQueue(name string, f func(map[string]interface{}, []byte) error) error
}

// client is a struct which holds all necessary data for RabbitMQ client
type client struct {
	c *amqp.Connection
}

func NewConnection(options *Options) (*amqp.Connection, error) {
	vh, _ := url.Parse(options.VirtualHost)
	dns := fmt.Sprintf("amqp://%s:%s@%s/%s", options.Username, options.Password, options.Addresses, vh)
	conn, err := amqp.Dial(dns)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Connect connects to RabbitMQ by dsn and return client object which uses openned client during function calls issued later in code
func NewClient(conn *amqp.Connection, options *Options) (Client, error) {
	c := &client{conn}
	if err := c.CreateScheme(options); err != nil {
		return nil, err
	}
	return c, nil
}

// CreateScheme creates all exchanges, queues and bindinges between them as specified in yaml string
func (c *client) CreateScheme(s *Options) error {
	ch, err := c.c.Channel()
	if err != nil {
		return err
	}

	// Create exchanges according to settings
	for name, e := range s.Exchanges {

		if e.Type == DelayedMessage {
			if _, ok := e.Args[DelayedType]; !ok {
				e.Args[DelayedType] = DefaultDelayType
			}
		}

		err = ch.ExchangeDeclarePassive(name, e.Type, e.Durable, e.AutoDelete, e.Internal, e.Nowait, e.Args)
		if err == nil {
			continue
		}
		ch, err = c.c.Channel()
		if err != nil {
			return err
		}

		err = ch.ExchangeDeclare(name, e.Type, e.Durable, e.AutoDelete, e.Internal, e.Nowait, e.Args)
		if err != nil {
			return err
		}
	}

	// Create queues according to settings
	for name, q := range s.Queues {
		_, err := ch.QueueDeclarePassive(name, q.Durable, q.AutoDelete, q.Exclusive, q.Nowait, nil)
		if err == nil {
			continue
		}

		ch, err = c.c.Channel()
		if err != nil {
			return err
		}

		_, err = ch.QueueDeclare(name, q.Durable, q.AutoDelete, q.Exclusive, q.Nowait, nil)
		if err != nil {
			return err
		}
	}

	// Create bindings only now that everything is setup.
	// (If the bindings were created in one run together with exchanges and queues,
	// it would be possible to create binding to not yet existent queue.
	// This way it's still possible but now is an error on the user side)
	for name, e := range s.Exchanges {
		for _, b := range e.Bindings {
			err = ch.ExchangeBind(name, b.Key, b.Exchange, b.Nowait, nil)
			if err != nil {
				return err
			}
		}
	}

	for name, q := range s.Queues {
		for _, b := range q.Bindings {
			err = ch.QueueBind(name, b.Key, b.Exchange, b.Nowait, nil)
			if err != nil {
				return err
			}
		}
	}

	ch.Close()
	return nil
}

// DeleteScheme deletes all queues and exchanges (together with bindings) as specified in yaml string
func (c *client) DeleteScheme(s *Options) error {
	ch, err := c.c.Channel()
	if err != nil {
		return err
	}

	for name := range s.Exchanges {
		err = ch.ExchangeDelete(name, false, false)
		if err != nil {
			return err
		}
	}

	for name := range s.Queues {
		_, err = ch.QueueDelete(name, false, false, false)
		if err != nil {
			return err
		}
	}
	ch.Close()
	return nil
}

// Close closes client to RabbitMQ
func (c *client) Close() error {
	return c.c.Close()
}

// SendMessage publishes plain text message to an exchange with specific routing key
func (c *client) SendMessage(ex, key string, data []byte, headers map[string]interface{}) error {
	ch, err := c.c.Channel()
	if err != nil {
		return err
	}

	err = ch.Publish(ex, key, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
			Headers:      headers,
		})
	if err != nil {
		return err
	}
	return ch.Close()
}

// SendBlob publishes byte blob message to an exchange with specific routing key
func (c *client) SendBlob(ex, key string, msg []byte) error {
	ch, err := c.c.Channel()
	if err != nil {
		return err
	}

	err = ch.Publish(ex, key, false, false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/octet-stream",
			Body:         msg,
		})
	if err != nil {
		return err
	}
	return ch.Close()
}

// ProcessQueue calls handler function on each message delivered to a queue
func (c *client) ProcessQueue(name string, f func(map[string]interface{}, []byte) error) error {
	ch, err := c.c.Channel()
	if err != nil {
		return err
	}
	msgs, err := ch.Consume(name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}
	for d := range msgs {
		err = f(d.Headers, d.Body)
		if err != nil {
			d.Nack(false, true)
			continue
		}
		d.Ack(false)
	}
	return nil
}
