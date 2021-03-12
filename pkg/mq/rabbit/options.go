package rabbit

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Exchange is structure with specification of properties of RabbitMQ exchange
type Exchange struct {
	Durable    bool      `yaml:"durable"`
	AutoDelete bool      `yaml:"auto-delete"`
	Internal   bool      `yaml:"internal"`
	Nowait     bool      `yaml:"nowait"`
	Type       string    `yaml:"type"`
	Bindings   []Binding `yaml:"bindings"`
	Args       map[string]interface{}
}

// Binding specifies to which exchange should be an exchange or a queue binded
type Binding struct {
	Exchange string `yaml:"exchange"`
	Key      string `yaml:"key"`
	Nowait   bool   `yaml:"nowait"`
}

// QueueSpec is a specification of properties of RabbitMQ queue
type QueueSpec struct {
	Durable    bool      `yaml:"durable"`
	AutoDelete bool      `yaml:"auto-delete"`
	Nowait     bool      `yaml:"nowait"`
	Exclusive  bool      `yaml:"exclusive"`
	Bindings   []Binding `yaml:"bindings"`
}

// Settings is a specification of all queues and exchanges together with all bindings.
type Options struct {
	Addresses   string               `yaml:"addresses"`
	Username    string               `yaml:"username"`
	Password    string               `yaml:"password"`
	VirtualHost string               `yaml:"virtual-host"`
	Exchanges   map[string]Exchange  `yaml:"exchanges"`
	Queues      map[string]QueueSpec `yaml:"queues"`
}

var defaultOptions = &Options{
	Addresses:   "127.0.0.1",
	Username:    "root",
	Password:    "root",
	VirtualHost: "/",
	Exchanges:   nil,
	Queues:      nil,
}

func NewOptions(v *viper.Viper) (*Options, error) {
	if err := v.UnmarshalKey("rabbit", defaultOptions); err != nil {
		return nil, errors.Wrap(err, "unmarshal rabbitmq option error")
	}
	return defaultOptions, nil
}
