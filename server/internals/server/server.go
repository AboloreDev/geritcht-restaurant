package server

import (
	"net/http"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Server struct {
	config   *config.Config
	database *gorm.DB
	logger   zerolog.Logger
}

func NewServer(
	config *config.Config,
	database *gorm.DB,
	logger zerolog.Logger) *Server {
	return &Server{
		config:   config,
		logger:   logger,
		database: database,
	}
}

func (s *Server) SetUpRoutes() *gin.Engine {
	router := gin.New()

	// Static Middlewares
	router.Use(s.CORS())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// ROUTES
	// Health check route
	router.GET("/health", s.HealthCheck)

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
