package di

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/storage/db"
)

func init() {
	app.RegisterProviders(db.NewOptions, db.NewDB, db.NewGormDB)
}
