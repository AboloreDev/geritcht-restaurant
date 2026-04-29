package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/database"
	"github.com/AboloreDev/geritcht-restaurant/internals/logger"
	"github.com/AboloreDev/geritcht-restaurant/internals/server"
	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("Geritcht Restaurant is starting")
}

func main() {
	// Initialise logger
	log := logger.New()

	// Load the env
	cfg, err := config.LoadEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("could not load config")
	}

	// Set Gin Mode
	gin.SetMode(cfg.Server.GinMode)

	// load the database configurations
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer mainDB.Close()

	// Initialise Server
	srv := server.NewServer(cfg, db, log)

	router := srv.SetUpRoutes()
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	// INITIATE GRACEFUL SHUTDOWN
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown server")
		return
	}

	log.Info().Msg("Shutting down server")
}
