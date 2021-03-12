package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"os"
	"sync"
)

type ProviderFactory interface {
	NewProvider() interface{}
}

type Middleware interface {
	NewHandler() gin.HandlerFunc
}

type Handler interface {
	MapRoutes(group *gin.RouterGroup)
}

type Server interface {
	Start() error
}

type HealthChecker interface {
}

type HealthServer interface {
	Start() error
}

type Application interface {
	Start() error
	RegisterProviders(providers ...interface{})
	RegisterMiddlewareProviders(providers ...interface{})
	RegisterHealthCheckProviders(providers ...interface{})
	RegisterHandlerProviders(providers ...interface{})
}

var app = NewApplication()

func NewApplication() Application {
	return &application{
		container:              dig.New(),
		providers:              make([]interface{}, 0, 32),
		handlerProviders:       make([]interface{}, 0, 16),
		healthCheckerProviders: make([]interface{}, 0, 8),
		middlewareProviders:    make([]interface{}, 0, 2),
		lock:                   sync.Mutex{},
	}
}

func Start() (err error) {
	return app.Start()
}

func RegisterProviders(providers ...interface{}) {
	app.RegisterProviders(providers...)
}

func RegisterMiddlewareProviders(providers ...interface{}) {
	app.RegisterMiddlewareProviders(providers...)
}

func RegisterHealthCheckProviders(providers ...interface{}) {
	app.RegisterHealthCheckProviders(providers...)
}

func RegisterHandlerProviders(providers ...interface{}) {
	app.RegisterHandlerProviders(providers...)
}

type application struct {
	container *dig.Container

	providers              []interface{}
	handlerProviders       []interface{}
	healthCheckerProviders []interface{}
	middlewareProviders    []interface{}

	lock sync.Mutex
}

func (a *application) Start() (err error) {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.buildContainer()

	var healthServer HealthServer
	if err := a.container.Invoke(func(s HealthServer) {
		healthServer = s
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to get health check server: %s", err)
		os.Exit(-1)
	}

	go func() {
		if err := healthServer.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to start health check server: %s", err)
			os.Exit(-1)
		}
	}()

	var svr Server
	if err := a.container.Invoke(func(f Server) {
		svr = f
	}); err != nil {
		fmt.Fprintf(os.Stderr, "failed to get factory: %s", err)
		os.Exit(-1)
	}

	return svr.Start()
}

func (a *application) RegisterProviders(providers ...interface{}) {
	a.providers = append(a.providers, providers...)
}

func (a *application) RegisterMiddlewareProviders(providers ...interface{}) {
	a.middlewareProviders = append(a.middlewareProviders, providers...)
}

func (a *application) RegisterHealthCheckProviders(providers ...interface{}) {
	a.healthCheckerProviders = append(a.healthCheckerProviders, providers...)
}

func (a *application) RegisterHandlerProviders(providers ...interface{}) {
	a.handlerProviders = append(a.handlerProviders, providers...)
}

func (a *application) buildContainer() error {
	if err := a.provide(a.providers); err != nil {
		return errors.Wrap(err, "failed to provide normal providers")
	}
	if err := a.provideWithGroup(a.middlewareProviders, "middleware"); err != nil {
		return errors.Wrap(err, "failed to provide middleware providers")
	}
	if err := a.provideWithGroup(a.handlerProviders, "handler"); err != nil {
		return errors.Wrap(err, "failed to provide handler providers")
	}
	if err := a.provideWithGroup(a.healthCheckerProviders, "health"); err != nil {
		return errors.Wrap(err, "failed to provide health providers")
	}
	return nil
}

func (a *application) provide(constructors []interface{}) error {
	for _, c := range constructors {
		if err := a.container.Provide(c); err != nil {
			return err
		}
	}
	return nil
}

func (a *application) provideWithGroup(constructors []interface{}, group string) error {
	for _, c := range constructors {
		if err := a.container.Provide(c, dig.Group(group)); err != nil {
			return err
		}
	}
	return nil
}
