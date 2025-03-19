package middleware

import (
	"net/http"
	"time"

	"fx-service/internal/reply"
	"fx-service/pkg/config"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// FiberRateLimiter limits the number of requests per IP address for Fiber router
func FiberRateLimiter(cfg config.RateLimiterConfig) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        cfg.MaxRequests,
		Expiration: time.Duration(cfg.Timeframe) * time.Second,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":  "Rate limit exceeded",
				"result": nil,
			})
		},
	})
}

// GinRateLimiter limits the number of requests per IP address for Gin router
func GinRateLimiter(cfg config.RateLimiterConfig) gin.HandlerFunc {
	rl := newGinRateLimiter()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		v := rl.getVisitor(ip, cfg.Timeframe)

		if v.count >= cfg.MaxRequests {
			payload := reply.Error("Rate limit exceeded")
			c.JSON(http.StatusTooManyRequests, payload)
			c.Abort()
			return
		}

		v.count++
		c.Next()
	}
}
