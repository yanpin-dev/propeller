package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/logger"
)

func init() {
	app.RegisterProviders(logger.NewOptions, logger.NewLogger)
}
