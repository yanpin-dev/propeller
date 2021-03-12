package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/storage/redis"
)

func init() {
	app.RegisterProviders(redis.NewOptions, redis.NewClient)
}
