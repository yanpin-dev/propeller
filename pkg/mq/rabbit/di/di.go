package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/mq/rabbit"
)

func init() {
	app.RegisterProviders(rabbit.NewOptions, rabbit.NewConnection, rabbit.NewClient)
}
