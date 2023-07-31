package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/arthureichelberger/logs/pkg/psql"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Packages
	db, err := psql.Connect(
		ctx,
		env("POSTGRES_USER", "logs"),
		env("POSTGRES_PASSWORD", "logs"),
		env("POSTGRES_HOST", "localhost"),
		env("POSTGRES_PORT", "5432"),
		env("POSTGRES_DATABASE", "logs"),
	)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return
	}
	_ = db

	router := gin.New()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", env("HTTP_PORT", "8080")),
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("failed to start http server")
		}
	}()

	<-done
	log.Debug().Msg("shutting down http server")
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown http server")
	}
}

func env(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
