package main

import (
	"context"
	"log"
	"time"

	"github.com/abstractbreazy/rate-limited-api/config"
	"github.com/abstractbreazy/rate-limited-api/internal/repository/redis"
	"github.com/abstractbreazy/rate-limited-api/pkg/server"
	"github.com/abstractbreazy/rate-limited-api/pkg/sigint"

	"github.com/abstractbreazy/rate-limited-api/internal/delivery/http"
)

func main() {
	// Init configurations
	conf, err := config.New().Init()
	if err != nil {
		log.Fatalf("%s", err.Error())
	}

	// Init Redis
	rd, err := redis.New(conf.Redis)
	if err != nil {
		log.Panic("Init: redis")
	}
	defer rd.Close()

	// Init http handler that implement endpoints
	handler := http.NewHandler(conf.RateLimiter, rd)

	// Create new echo instances with implemented endpoins.
	routes := http.NewRoutes(handler)

	srv := server.New(conf.Server, routes)
	if err != nil {
		log.Panic("creating HTTP server", err)
	}

	// Run server in goroutine.
	srv.Run()

	log.Println("service is started:", conf.Server.Port)

	sigint.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("stopping service ...")

	// Shutdown Server.
	err = srv.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
