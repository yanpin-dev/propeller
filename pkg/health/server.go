package health

import (
	"github.com/yanpin-dev/propeller/pkg/app"
	"github.com/yanpin-dev/propeller/pkg/logger"
	"fmt"
	"net/http"
	"time"
)

func NewServer(options *Options, log logger.LogInfoFormat, checkers []Checker) app.HealthServer {
	return &server{
		log:      log,
		host:     "0.0.0.0",
		port:     8081,
		checkers: checkers,
	}
}

type server struct {
	log      logger.LogInfoFormat
	host     string
	port     uint16
	checkers []Checker

	httpServer *http.Server
}

func (s *server) Start() error {
	handler := NewHandler()
	for _, c := range s.checkers {
		handler.AddChecker(c)
	}
	http.Handle("/health/", handler)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	s.log.Infof("health check server start on %s", addr)

	s.httpServer = &http.Server{
		Addr:        addr,
		Handler:     http.DefaultServeMux,
		ReadTimeout: time.Duration(10) * time.Second,
	}
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
