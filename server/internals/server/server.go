package server

import (
	"net/http"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/interfaces"
	"github.com/AboloreDev/geritcht-restaurant/internals/services"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	cfg                 *config.Config
	log                 zerolog.Logger
	db                  *gorm.DB
	authServices        *services.AuthService
	redisStore          interfaces.Cacher
	uploadServices      *services.UploadService
	categoryServices    *services.CategoryService
	menuServices        *services.MenuService
	allergenServices    *services.AllergenService
	dietaryTagsService  *services.DietaryTagsService
	userServices        *services.UserService
	reservationServices *services.ReservationService
	waitlistService     *services.WaitlistService
	tableServices       *services.TableService
	paymentService      *services.PaymentService
	orderService        *services.OrderService
}

func NewServer(
	cfg *config.Config,
	db *gorm.DB,
	log zerolog.Logger,
	authServices *services.AuthService,
	redisStore interfaces.Cacher,
	uploadServices *services.UploadService,
	categoryServices *services.CategoryService,
	menuServices *services.MenuService,
	allergenServices *services.AllergenService,
	dietaryTagsService *services.DietaryTagsService,
	reservationServices *services.ReservationService,
	waitlistService *services.WaitlistService,
	userServices *services.UserService,
	tableServices *services.TableService,
	paymentService *services.PaymentService,
	orderService *services.OrderService) *Server {
	return &Server{
		cfg:                 cfg,
		log:                 log,
		db:                  db,
		authServices:        authServices,
		redisStore:          redisStore,
		uploadServices:      uploadServices,
		categoryServices:    categoryServices,
		menuServices:        menuServices,
		allergenServices:    allergenServices,
		dietaryTagsService:  dietaryTagsService,
		userServices:        userServices,
		waitlistService:     waitlistService,
		reservationServices: reservationServices,
		tableServices:       tableServices,
		paymentService:      paymentService,
		orderService:        orderService,
	}
}

func (s *Server) SetUpRoutes() *gin.Engine {
	router := gin.New()

	// Static Middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(s.CORS())

	// ROUTES
	// Health check route
	router.GET("/health", s.HealthCheck)

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			// Auth Routes
			auth.POST("/register", s.RegisterUserHandler)
			auth.POST("/login", s.LoginUserHandler)
			auth.POST("/logout", s.LogoutHandler)
			auth.POST("/refresh", s.RefreshTokenHandler)
			auth.POST("/forgot", s.ForgotPasswordHandler)
			auth.POST("/reset-password", s.ResetPasswordHandler)
			auth.POST("/verify", s.VerifyEmailHandler)
			auth.POST("/verify-reset-otp", s.VerifyResetOTPHandler)

		}
		protected := api.Group("/")
		protected.Use(s.AuthMiddleware())
		{
			users := protected.Group("/users")
			{
				// User Protected Routes
				users.PATCH("/password-change", s.ChangePasswordHandler)
				users.GET("/profile", s.GetUserProfileHandler)
				users.PATCH("/profile", s.UpdateUserProfileHandler)
				users.PATCH("/profile/deactivate", s.DeactivateUserHandler)
				users.PATCH("/profile/activate", s.ActivateUserHandler)
				users.GET("/", s.AdminMiddleware(), s.GetAllUserHandler)
			}

			staffs := protected.Group("/staff")
			{
				// Staff Protected Routes
				staffs.GET("/profile", s.StaffMiddleware(), s.GetStaffProfileHandler)
				staffs.PATCH("/profile", s.StaffMiddleware(), s.UpdateStaffProfileHandler)
				staffs.PATCH("/profile/deactivate", s.StaffMiddleware(), s.DeactivateStaffHandler)
				staffs.PATCH("/profile/activate", s.StaffMiddleware(), s.ActivateStaffHandler)
				staffs.GET("/", s.AdminMiddleware(), s.GetAllStaffsHandler)
			}

			category := protected.Group("/categories")
			{
				// Category Protected Routes
				category.POST("/", s.AdminMiddleware(), s.CreateCategoryHandler)
				category.PATCH("/:id", s.AdminMiddleware(), s.UpdateCategoryHandler)
				category.DELETE("/:id", s.AdminMiddleware(), s.DeleteCategory)
				category.GET("/categories", s.GetCategoriesHandler)
				category.GET("/categories/:id", s.GetCategory)
			}

			menu := protected.Group("/menu")
			{
				// Menu Protected Routes
				menu.POST("/", s.AdminMiddleware(), s.CreateMenuHandler)
				menu.PATCH("/:id", s.AdminMiddleware(), s.UpdateMenuHandler)
				menu.GET("/menu", s.GetAllMenuHandler)
				menu.GET("/menu/:id", s.GetMenuHandler)
				menu.DELETE("/:id", s.AdminMiddleware(), s.DeleteMenuHandler)
				menu.PATCH("/:id/toggle", s.AdminMiddleware(), s.ToggleMenuAvailabilityHandler)
				menu.POST("/:id/images", s.AdminMiddleware(), s.UploadMenuImageHandler)
				menu.DELETE("/images/:id", s.AdminMiddleware(), s.DeleteMenuImageHandler)
			}

			allergens := protected.Group("/allergens")
			{
				// Allergens and Dietary Tags Protected Routes
				allergens.POST("/", s.AdminMiddleware(), s.CreateAllergenHandler)
				allergens.GET("/", s.AdminMiddleware(), s.GetAllAllergenHandler)
				allergens.PATCH("/:id", s.AdminMiddleware(), s.UpdateAllegenHandler)
				allergens.DELETE("/:id", s.AdminMiddleware(), s.DeleteAllergenHandler)

			}

			tags := protected.Group("/tags")
			{
				tags.POST("/", s.AdminMiddleware(), s.CreateDietaryTagHandler)
				tags.GET("/", s.AdminMiddleware(), s.GetAllDietaryTagHandler)
				tags.PATCH("/:id", s.AdminMiddleware(), s.UpdateDietaryTagHandler)
				tags.DELETE("/:id", s.AdminMiddleware(), s.DeleteDietaryTagHandler)
			}

			table := protected.Group("/tables")
			{
				// Table Protected Routes
				table.POST("/", s.AdminMiddleware(), s.CreateTableHandler)
				table.GET("/", s.GetAllTablesHandler)
				table.GET("/:id", s.GetTableHandler)
				table.PATCH("/:id", s.AdminMiddleware(), s.UpdateTableHandler)
				table.DELETE("/:id", s.AdminMiddleware(), s.DeleteTableHandler)
			}

			reservation := protected.Group("/reservations")
			{
				// Reservation Protected Routes
				reservation.POST("/", s.CreateReservationHandler)
				reservation.GET("/", s.RoleMiddleware("admin", "staff"), s.GetAllReservationsHandler)
				reservation.GET("/:id/user", s.GetUserReservationHandler)
				reservation.GET("/user", s.GetAllUserReservationsHandler)
				reservation.GET("/today", s.GetTodayReservationHandler)
				reservation.POST("/:id", s.CheckInReservationHandler)
				reservation.POST("/:id/cancel", s.RoleMiddleware("admin", "staff"), s.CancelReservationHandler)
				reservation.GET("/availability", s.CheckAvailabilityHandler)
			}
		}

		// Public routes
		api.GET("/categories", s.GetCategoriesHandler)
		api.GET("/categories/:id", s.GetCategory)
		api.GET("/menu", s.GetAllMenuHandler)
		api.GET("/menu/:id", s.GetMenuHandler)
		api.GET("/table", s.GetAllTablesHandler)
		api.GET("/table/:id", s.GetTableHandler)
		api.GET("/availability", s.CheckAvailabilityHandler)
	}

	return router
}

func (s *Server) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *Server) CORS() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}
