package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/health"
	"github.com/yanpin-dev/propeller/pkg/health/checker"
	"github.com/yanpin-dev/propeller/pkg/logger"
	"go.uber.org/dig"
)

func init() {
	app.RegisterHealthCheckProviders(
		checker.NewMySQLChecker,
		checker.NewDefaultRedisChecker,
		checker.NewRabbitChecker,
	)
	app.RegisterProviders(
		NewServer,
		health.NewOptions,
	)
}

type Group struct {
	dig.In

	Options  *health.Options
	Logger   logger.LogInfoFormat
	Checkers []health.Checker `group:"health"`
}

func NewServer(g Group) app.HealthServer {
	return health.NewServer(
		g.Options,
		g.Logger,
		g.Checkers,
	)
}
