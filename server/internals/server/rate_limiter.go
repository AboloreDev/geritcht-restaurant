package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

// @Header 200 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 200 {string} X-RateLimit-Remaining "Remaining requests in the current window"
// @Header 429 {string} X-RateLimit-Limit "Maximum requests allowed"
// @Header 429 {string} X-RateLimit-Remaining "Remaining requests (0 when limited)"
// Under limit -> accept/ allow
func (s *Server) RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		key := fmt.Sprintf("rate_limit:%s:%s", ctx.FullPath(), ip)

		count, err := s.redisStore.Increment(ctx, key)
		if err != nil {
			// Allow access if redis is unavailable
			ctx.Next()
			return
		}

		if count == 1 {
			// Start the redis count from the first request that comes in
			s.redisStore.Expire(ctx, key, window)
		}

		// Over limit - reject
		if int(count) > limit {
			ctx.Header("X-RateLimit-Limit", strconv.Itoa(limit))
			ctx.Header("X-RateLimit-Remaining", "0")
			utils.TooManyRequests(ctx, fmt.Sprintf("Rate limit exceeded. Maximum %d requests every %s.", limit, window), nil)
			ctx.Abort()
			return
		}

		
		ctx.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(limit-int(count)))
		ctx.Next()
	}
}
