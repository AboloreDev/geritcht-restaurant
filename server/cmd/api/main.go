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
	"github.com/AboloreDev/geritcht-restaurant/internals/publisher"
	redisImport "github.com/AboloreDev/geritcht-restaurant/internals/redis"
	"github.com/AboloreDev/geritcht-restaurant/internals/server"
	"github.com/AboloreDev/geritcht-restaurant/internals/services"
	websockets "github.com/AboloreDev/geritcht-restaurant/internals/web-sockets"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func init() {
	fmt.Println("Geritcht Restaurant is starting")
}

func main() {
	// Initialise logger
	log := logger.New()
	// App Context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Worker context
	workerCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	opt, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid Redis URL")
	}
	client := redis.NewClient(opt)

	// Initialise Redis
	var redisStore interfaces.Cacher
	if cfg.Redis.URL == "" {
		log.Warn().Msg("No Redis URL — using NopStore")
		redisStore = redisImport.NewNopCache()
	} else {
		if err := client.Ping(ctx).Err(); err != nil {
			log.Warn().Err(err).Msg("Redis ping failed — using NopStore")
			redisStore = redisImport.NewNopCache()
		} else {
			log.Info().Msg("Redis connected successfully")
			redisStore = redisImport.NewStore(client)
		}
	}

	defer func() {
		if err := client.Close(); err != nil {
			log.Error().Err(err).Msg("error closing redis connection")
		}
	}()

	// Upload Provider
	var uploadProvider interfaces.UploadProvider
	if cfg.Uploads.UploadProvider == "cloudinary" {
		uploadProvider = providers.NewCloudinaryUploader(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Uploads.UploadPath)
	}

	// Event Publisher Initialisation
	eventPublisher, err := publisher.NewEventPublisher(
		ctx,
		&config.RedisConfig{
			QUEUE_NAME: cfg.Redis.QUEUE_NAME,
		}, client)
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialise events")
		return
	}
	defer eventPublisher.CloseMessage()

	// Services
	authServices := services.NewAuthService(db, cfg, eventPublisher)
	uploadServices := services.NewUploadServices(uploadProvider)
	categoryServices := services.NewCategoryService(db, redisStore)
	menuServices := services.NewMenuService(db, redisStore)
	allegenServices := services.NewAllergenService(db, redisStore)
	dietaryTagsService := services.NewDietaryTagsService(db, redisStore)
	userServices := services.NewUserService(db)
	reservationServices := services.NewReservationService(db, redisStore, eventPublisher)
	waitlistServices := services.NewWaitlistService(db)
	tableServices := services.NewTableService(db, redisStore)
	paymentService := services.NewPaymentService(db, redisStore, eventPublisher, &config.Config{
		Paystack: config.PaystackConfig{
			PaystackSecretKey: cfg.Paystack.PaystackSecretKey,
			PaystackPublicKey: cfg.Paystack.PaystackPublicKey,
		},
	})
	orderService := services.NewOrderService(db, redisStore)
	cartService := services.NewCartService(db)
	// websockts hub for order
	hub := websockets.NewHub()

	// DB Workers
	noShowWorker := services.NewNoShowWorker(db, eventPublisher, redisStore)
	reminderWorker := services.NewReminderWorker(db, redisStore, eventPublisher)
	outboxWorker := services.NewOutboxWorker(db, eventPublisher)
	orderStatusAutoWorker := services.NewOrderAutoWorker(db, hub)

	// Go routines for worker
	go noShowWorker.StartMarkNoShowWorker(workerCtx, log)
	go reminderWorker.StartReminderWorker(workerCtx, log)
	go outboxWorker.StartOutboxWorker(workerCtx, log)
	go orderStatusAutoWorker.StartOrderUpdateWorker(workerCtx, log)

	
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
		reservationServices,
		waitlistServices,
		userServices,
		tableServices,
		paymentService,
		orderService,
		cartService,
		hub,
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
	defer cancel()

	err = httpServer.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to shutdown server")
		return
	}

	log.Info().Msg("Server shutting down")
}
