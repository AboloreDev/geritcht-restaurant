package server

import (
	"strings"

	"github.com/AboloreDev/geritcht-restaurant/internals/models"
	"github.com/AboloreDev/geritcht-restaurant/internals/utils"

	"github.com/gin-gonic/gin"
)

// @Security BearerAuth
// @Failure 401 {object} utils.Response "Unauthorized"
func (s *Server) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Authorization Bearer token
		authHeader := ctx.GetHeader("authorization")
		if authHeader == "" {
			utils.UnAuthorized(ctx, "Authorization header required", nil)
			ctx.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			utils.UnAuthorized(ctx, "Invalid Header Format", nil)
			ctx.Abort()
			return
		}

		claims, err := utils.ValidateToken(tokenParts[1], s.cfg.JWT.JWTSecret)
		if err != nil {
			utils.UnAuthorized(ctx, "Invalid Token", err)
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Set("user_role", claims.Role)

		ctx.Next()
	}
}

// @Security BearerAuth
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
func (s *Server) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")

		if !exists {
			utils.Forbidden(ctx, "Forbidden", nil)
			ctx.Abort()
			return
		}

		if role != string(models.RoleAdmin) {
			utils.Forbidden(ctx, "Forbidden", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// @Security BearerAuth
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
func (s *Server) StaffMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("user_role")

		if !exists {
			utils.Forbidden(ctx, "Forbidden", nil)
			ctx.Abort()
			return
		}

		if role != string(models.RoleStaff) {
			utils.Forbidden(ctx, "Forbidden", nil)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// @Security BearerAuth
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 403 {object} utils.Response "Forbidden"
func (s *Server) RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if len(roles) == 0 {
			ctx.Next()
			return
		}

		userRole, exists := ctx.Get("user_role")
		if !exists {
			utils.Forbidden(ctx, "You do not have permission to access this resource", nil)
			ctx.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.Forbidden(ctx, "Invalid role information", nil)
			ctx.Abort()
			return
		}

		for _, role := range roles {
			if roleStr == role {
				ctx.Next()
				return
			}
		}

		utils.Forbidden(ctx, "Forbidden", nil)
		ctx.Abort()
	}
}
