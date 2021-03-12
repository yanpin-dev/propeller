package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/logger"
	"github.com/yanpin-dev/propeller/pkg/server"
	"go.uber.org/dig"
)

func init() {
	app.RegisterProviders(server.NewOptions, NewServer)
}

type FactoryGroup struct {
	dig.In

	Options     *server.Options
	Logger      logger.LogInfoFormat
	Handlers    []app.Handler    `group:"handler"`
	Middlewares []app.Middleware `group:"middleware"`
}

func NewServer(g FactoryGroup) app.Server {
	return server.NewServer(
		g.Options,
		g.Logger,
		g.Handlers,
		g.Middlewares,
	)
}
