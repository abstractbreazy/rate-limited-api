package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Config struct {
	Port         string        `yaml:"port" mapstructure:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout" mapstructure:"write_timeout"`
}

func NewConfig() *Config {
	return &Config{}
}

type Server struct {
	echo *echo.Echo
}

func New(conf *Config, router http.Handler) (s *Server) {
	s = new(Server)
	s.echo = echo.New()
	s.echo.Server = &http.Server{
		Addr:         ":" + conf.Port,
		Handler:      router,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
	}
	return
}

// Run the server in goroutine.
func (s *Server) Run() {
	go s.run()
}

// Shutdown entire server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}

func (s *Server) Close() (err error) {
	return s.echo.Close()
}

func (s *Server) run() {
	if err := s.echo.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}
}
