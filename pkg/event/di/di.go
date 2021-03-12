package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/event"
)

func init() {
	app.RegisterProviders(event.NewOptions, event.NewEventPublisher)
}
