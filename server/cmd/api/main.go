package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/database"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/logger"
	"github.com/AboloreDev/geritcht-restaurant/internals/providers"
	redisImport "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/AboloreDev/geritcht-restaurant/internals/server"
	"github.com/AboloreDev/geritcht-restaurant/internals/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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

	// // load the database configurations
	db, err := database.New(&cfg.Database)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to database")
	}

	mainDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer mainDB.Close()

	// Upload Provider
	var uploadProvider interfaces.UploadProvider
	if cfg.Uploads.UploadProvider == "cloudinary" {
		uploadProvider = providers.NewCloudinaryUploader(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Uploads.UploadPath)
	}

	// Initialise Redis
	var redisStore interfaces.Cacher
	if cfg.Redis.URL == "" {
		log.Warn().Msg("No Redis URL — using NopStore")
		redisStore = redisImport.NewNopCache()
	} else {
		client := redis.NewClient(&redis.Options{
			Addr: cfg.Redis.URL,
			DB:   0,
		})

		redisCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := client.Ping(redisCtx).Err(); err != nil {
			log.Warn().Err(err).Msg("Redis ping failed — using NopStore")
			redisStore = redisImport.NewNopCache()
		} else {
			log.Info().Msg("Redis connected successfully")
			redisStore = redisImport.NewStore(client)
		}
	}

	// Services
	authServices := services.NewAuthService(db, cfg)
	uploadServices := services.NewUploadServices(uploadProvider)
	categoryServices := services.NewCategoryService(db, redisStore)
	menuServices := services.NewMenuService(db, redisStore)
	allegenServices := services.NewAllergenService(db)
	dietaryTagsService := services.NewDietaryTagsService(db)
	userServices := services.NewUserService(db)

	// Initialise Server
	srv := server.NewServer(
		cfg, db, log,
		authServices,
		redisStore,
		uploadServices,
		categoryServices,
		menuServices,
		allegenServices,
		dietaryTagsService,
		userServices,
	)

	router := srv.SetUpRoutes()
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	// Start server
	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("http server started")
		err = httpServer.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("Failed to start http server")
		}
	}()

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

	log.Info().Msg("Server shutting down")
}
