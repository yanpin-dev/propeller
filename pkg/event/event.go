package event

import (
	"github.com/yanpin-dev/propeller/pkg/logger"
	xrabbit "github.com/yanpin-dev/propeller/pkg/mq/rabbit"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

const (
	ID        = "x-fcyp-id"
	Timestamp = "x-fcyp-timestamp"
	Type      = "x-fcyp-type"
	Source    = "x-fcyp-source"

	RoutingKey = "#"
)

type DomainEvent interface {
	GetEventID() string
	GetAggregateID() string
	GetWhen() *time.Time
	GetType() string
}
type DefaultDomainEvent struct {
	EventID     string     `json:"-"`
	AggregateID string     `json:"aggregateId"`
	When        *time.Time `json:"-"`
	Type        string     `json:"-"`
}

func (d *DefaultDomainEvent) GetEventID() string {
	return d.EventID
}
func (d *DefaultDomainEvent) GetAggregateID() string {
	return d.AggregateID
}
func (d *DefaultDomainEvent) GetWhen() *time.Time {
	return d.When
}
func (d *DefaultDomainEvent) GetType() string {
	return d.Type
}

func NewEventPublisher(log logger.LogInfoFormat, c xrabbit.Client, options *Options) Publisher {
	return &RabbitEventPublisher{
		log:        log,
		connection: c,
		source:     options.Source,
		exchange:   options.Exchange,
	}
}

type Publisher interface {
	Publish(interface{}) error
}

type RabbitEventPublisher struct {
	log        logger.LogInfoFormat
	connection xrabbit.Client
	source     string
	exchange   string
}

func (r *RabbitEventPublisher) Publish(event interface{}) error {
	headers := make(map[string]interface{})
	if evt, ok := event.(DomainEvent); ok {
		eID := evt.GetEventID()
		if evt.GetEventID() == "" {
			eID = uuid.New().String()
		}

		t := evt.GetWhen()
		if evt.GetWhen() == nil {
			now := time.Now()
			t = &now
		}
		headers[ID] = eID
		headers[Type] = evt.GetType()
		headers[Timestamp] = t.UnixNano() / 1000000
		headers[Source] = r.source
	}
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = r.connection.SendMessage(r.exchange, RoutingKey, data, headers)

	hjson, _ := json.Marshal(headers)

	if err != nil {
		r.log.Warnf("publish event failed, header=%s, event=%s, error=%s", string(hjson), string(data), err.Error())
		return err
	}

	r.log.Debugf("publish event success, header=%v, event=%s", headers, string(data))
	return nil
}
