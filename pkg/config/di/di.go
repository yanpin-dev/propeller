package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/config"
)

func init() {
	app.RegisterProviders(config.NewViper)
}
