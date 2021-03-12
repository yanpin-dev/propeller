package event

import (
	"github.com/yanpin-dev/propeller/pkg/mq/rabbit"
	"errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
	"time"
)

func TestPublish(t *testing.T) {
	log, _ := NewLogger()
	p := &RabbitEventPublisher{
		log:        log,
		connection: &MockConnection{},
		source:     "",
		exchange:   "",
	}
	t2 := time.Now()
	event := &MockDomainEvent{
		DefaultDomainEvent: DefaultDomainEvent{
			EventID:     "event-id",
			AggregateID: "aggregate-id",
			When:        &t2,
			Type:        "test",
		},
		Data: "data",
	}
	p.Publish(event)
}

func NewLogger() (*zap.SugaredLogger, error) {

	cfg := zap.NewDevelopmentConfig()

	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	//cfg.OutputPaths = []string{o.Logger.FileName}
	log, err := cfg.Build()
	if err != nil {
		return nil, errors.New("zap logger build constructs failed.")
	}
	return log.Sugar(), nil
}

type MockConnection struct {
}

func (m *MockConnection) CreateScheme(s *rabbit.Options) error {
	return nil
}

func (m *MockConnection) ProcessQueue(name string, f func(map[string]interface{}, []byte) error) error {
	return nil
}

func (m *MockConnection) SendMessage(ex, key string, data []byte, headers map[string]interface{}) error {
	return nil
}

type MockDomainEvent struct {
	DefaultDomainEvent
	Data string `json:"data"`
}
