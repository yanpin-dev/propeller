package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/http/middleware"
)

func init() {
	app.RegisterMiddlewareProviders(
		middleware.NewErrorMiddleware,
	)
}
