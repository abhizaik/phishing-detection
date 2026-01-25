package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/abhizaik/SafeSurf/internal/service/cache"
	"github.com/gin-gonic/gin"
)

// RateLimiter returns a Gin middleware that limits requests per IP
func RateLimiter(limit int64, window time.Duration) gin.HandlerFunc {
	c, err := cache.New()
	if err != nil {
		fmt.Printf("Warning: Rate limiter failed to connect to cache: %v. Rate limiting disabled.\n", err)
		return func(ctx *gin.Context) { ctx.Next() }
	}

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		key := fmt.Sprintf("ratelimit:%s", ip)

		count, err := c.Increment(ctx, key)
		if err != nil {
			// On cache error, fail open to avoid blocking users
			ctx.Next()
			return
		}

		// If this is the first request in the window, set the expiration
		if count == 1 {
			_ = c.Expire(ctx, key, window)
		}

		// Set rate limit headers
		ctx.Header("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
		ctx.Header("X-RateLimit-Remaining", strconv.FormatInt(max(0, limit-count), 10))
		ctx.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))

		if count > limit {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests",
				"retry_after": window.Seconds(),
			})
			return
		}

		ctx.Next()
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
