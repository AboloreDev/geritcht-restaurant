package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AboloreDev/geritcht-restaurant/internals/utils"
	"github.com/gin-gonic/gin"
)

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
			utils.TooManyRequests(ctx, "Too many requests", err)
			ctx.Abort()
			return
		}

		// Under limit -> accept/ allow
		ctx.Header("X-RateLimit-Limit", strconv.Itoa(limit))
		ctx.Header("X-RateLimit-Remaining", strconv.Itoa(limit-int(count)))
		ctx.Next()
	}
}
