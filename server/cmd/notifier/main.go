package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/AboloreDev/geritcht-restaurant/cmd/notifier/subscriber"
	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/email"
	"github.com/AboloreDev/geritcht-restaurant/internals/logger"
	"github.com/redis/go-redis/v9"
)

func init() {
	log.Println("Starting notification micro-service")
}

func main() {
	log := logger.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load All config")
	}

	emailClient := email.NewResendEmailClient(
		ctx,
		&cfg.Resend,
	)

	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid Redis URL")
	}

	client := redis.NewClient(opt)

	defer func() {
		if err := client.Close(); err != nil {
			log.Error().Err(err).Msg("error closing redis connection")
		}
	}()

	sub, err := subscriber.NewEventSubscriber(
		ctx,
		&config.RedisConfig{QUEUE_NAME: cfg.Redis.QUEUE_NAME},
		client,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create subscriber")
	}

	// start consuming
	if err := sub.Start(ctx, emailClient); err != nil {
		log.Fatal().Err(err).Msg("failed to start subscriber")
	}

	// dummy HTTP server just to satisfy Render's Web Service port requirement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	log.Info().Msgf("notifier health server listening on :%s", port)
	http.ListenAndServe(":"+port, nil)

	log.Info().Msg("notification micro-service started")

	// Initiate graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("notification micro-service has started and waiting for signal")
	cancel()
}
