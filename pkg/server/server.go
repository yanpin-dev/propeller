package server

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
)

type server struct {
	engine      *gin.Engine
	options     *Options
	logger      logger.LogInfoFormat
	handlers    []app.Handler
	middlewares []app.Middleware
}

func NewServer(options *Options, log logger.LogInfoFormat, handlers []app.Handler, middlewares []app.Middleware) app.Server {
	e := gin.Default()
	e.RedirectTrailingSlash = true
	e.RedirectFixedPath = true
	e.HandleMethodNotAllowed = true
	return &server{
		engine:      e,
		options:     options,
		logger:      log,
		handlers:    handlers,
		middlewares: middlewares,
	}
}

func (s *server) Start() error {
	s.setupMiddlewares()
	s.setupRoutes()
	return s.engine.Run(fmt.Sprintf("%s:%d", s.options.Host, s.options.Port))
}

func (s *server) setupRoutes() {
	api := s.engine.Group("/api")
	for _, h := range s.handlers {
		h.MapRoutes(api)
	}
}

func (s *server) setupMiddlewares() {
	for _, m := range s.middlewares {
		s.engine.Use(m.NewHandler())
	}
}
