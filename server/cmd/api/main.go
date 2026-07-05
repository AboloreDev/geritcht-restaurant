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
	"github.com/AboloreDev/geritcht-restaurant/internals/repositories"
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

	// repositories
	authRepo := repositories.NewAuthRepository(db)
	userRepo := repositories.NewUserRepository(db)
	cartRepo := repositories.NewCartRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	tagRepo := repositories.NewDietaryTagRepository(db)
	allergenRepo := repositories.NewAllergenRepository(db)
	menuRepo := repositories.NewMenuRepository(db)
	orderRepo := repositories.NewOrderRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	tableRepo := repositories.NewTableRepository(db)
	reservationRepo := repositories.NewReservationRepository(db)
	waitlistRepo := repositories.NewWaitlistRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)
	ingredientRepo := repositories.NewIngredientRepository(db)
	outboxRepo := repositories.NewOutboxRepository(db)
	noShowWorkerRepo := repositories.NewReservationNoShowRepository(db)
	reminderRepo := repositories.NewReservationReminderRepository(db)

	// Services
	authServices := services.NewAuthService(cfg, eventPublisher, userRepo, authRepo, cartRepo)
	uploadServices := services.NewUploadServices(uploadProvider)
	categoryServices := services.NewCategoryService(redisStore, categoryRepo)
	menuServices := services.NewMenuService(menuRepo, redisStore)
	allegenServices := services.NewAllergenService(allergenRepo, redisStore)
	dietaryTagsService := services.NewDietaryTagsService(tagRepo, redisStore)
	userServices := services.NewUserService(userRepo)
	reservationServices := services.NewReservationService(db, reservationRepo, redisStore, eventPublisher)
	waitlistServices := services.NewWaitlistService(waitlistRepo)
	tableServices := services.NewTableService(tableRepo, redisStore)
	inventoryService := services.NewInventoryService(db, eventPublisher, redisStore, inventoryRepo)
	paymentService := services.NewPaymentService(db, redisStore, eventPublisher, &config.Config{
		Paystack: config.PaystackConfig{
			PaystackSecretKey: cfg.Paystack.PaystackSecretKey,
			PaystackPublicKey: cfg.Paystack.PaystackPublicKey,
		},
	}, inventoryRepo, &http.Client{Timeout: 10 * time.Second}, paymentRepo, *inventoryService)
	orderService := services.NewOrderService(db, orderRepo, paymentRepo, cartRepo, redisStore)
	cartService := services.NewCartService(cartRepo, menuRepo)
	ingredientService := services.NewIngredientService(redisStore, eventPublisher, ingredientRepo, userRepo, paymentRepo)
	recipesService := services.NewMenuItemIngredientService(db)

	// websockts hub for order
	hub := websockets.NewHub()

	// DB Workers
	noShowWorker := services.NewNoShowWorker(eventPublisher, redisStore, noShowWorkerRepo)
	reminderWorker := services.NewReminderWorker(reminderRepo, redisStore, eventPublisher)
	outboxWorker := services.NewOutboxWorker(outboxRepo, eventPublisher)
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
		ingredientService,
		recipesService,
		inventoryService,
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
