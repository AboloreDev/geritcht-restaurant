package server

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/config"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

var testConfig *config.Config
var testServer *Server

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	testConfig = &config.Config{
		JWT: config.JWTConfig{
			JWTSecret:                 "test_key",
			JWTTokenExpiration:        15 * time.Minute,
			JWTRefreshTokenExpiration: 7 * 24 * time.Hour,
		},
	}

	testServer = &Server{
		cfg: testConfig,
	}

	code := m.Run()
	os.Exit(code)
}

func SetupAuthTestRouter(s *Server) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	protected := router.Group("/")
	protected.Use(s.AuthMiddleware())
	{
		protected.GET("/test", func(ctx *gin.Context) {
			userID := ctx.GetUint("user_id")
			role := ctx.GetUint("user_role")
			ctx.JSON(http.StatusOK, gin.H{
				"user_id":   userID,
				"user_role": role,
			})
		})
	}

	return router
}

func SetupRoleTestRouter(s *Server, role string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	protected := router.Group("/")
	protected.Use(
		s.AuthMiddleware(),
		s.RoleMiddleware(role))
	{
		protected.GET("/test", func(ctx *gin.Context) {
			ctx.Get("user_role")
			ctx.JSON(http.StatusOK, gin.H{"message": "access granted"})
		})
	}

	return router
}

func Test_AuthMiddleware(t *testing.T) {
	router := SetupAuthTestRouter(testServer)
	email := "test@gmail.com"
	role := "user"

	validToken, _, _ := utils.GenerateTokenPair(&testConfig.JWT, 1, email, role)

	tests := []struct {
		name           string
		token          string
		expectedCode   int
		expectedUserID string
	}{
		{
			name:         "Valid token",
			token:        validToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "Invalid token",
			token:        "invalid_token",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Empty token",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "No token provided",
			token:        "Bearer ",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Expired token",
			token:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2Nzc0MDY3MzV9.1234567890",
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "Mising bearer prefix",
			token:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE2Nzc0MDY3MzV9.1234567890",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedUserID != "" {
				if !strings.Contains(w.Body.String(), tt.expectedUserID) {
					t.Errorf("expected response body to contain %q, got %q", tt.expectedUserID, w.Body.String())
				}
			}
		})
	}
}

func Test_RoleMiddleware(t *testing.T) {
	email := "admin@gmail.com"
	role := "admin"
	router := SetupRoleTestRouter(testServer, role)

	adminToken, _, _ := utils.GenerateTokenPair(
		&testConfig.JWT,
		1,
		email,
		role,
	)

	userToken, _, _ := utils.GenerateTokenPair(
		&testConfig.JWT,
		2,
		"user@gmail.com",
		"user",
	)

	tests := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{
			name:         "Valid admin token",
			token:        adminToken,
			expectedCode: http.StatusOK,
		},
		{
			name:         "User token accessing admin route",
			token:        userToken,
			expectedCode: http.StatusForbidden,
		},
		{
			name:         "No token provided",
			token:        "",
			expectedCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			req.Header.Set("Authorization", "Bearer "+tt.token)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, w.Code)
			}
		})
	}
}
