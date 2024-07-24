package middleware

import (
	"fmt"
	"fx-service/internal/service/stats"
	"fx-service/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// makeContextLogger helper used by all routers' middleware to create a new logger with the path and IP
func makeContextLogger(baseLogger *logger.Logger, path, ip string) *logger.Logger {
	fields := map[string]interface{}{
		"path": path,
		"ip":   ip,
	}
	return baseLogger.WithContextFields(fields)
}

// updateHitCount helper used by all routers' middleware to increment the global hit counter
func updateHitCount(path string) {
	st := stats.GetInstance()
	st.IncHitCount()
	st.IncPath(path)
}

// FiberLogger attaches a new contextual logger to the context, for Fiber router
func FiberLogger(baseLogger *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()

		// Create a new logger with the path and IP
		ctxLogger := makeContextLogger(baseLogger, path, c.IP())

		// Don't log calls to favicon
		if path != "/favicon.ico" {
			ctxLogger.Info(fmt.Sprintf("req: %s %s", c.Method(), c.OriginalURL()), nil)
		}

		// Set the logger in the context
		c.Locals("logger", ctxLogger)

		// Increment the global hit counter
		updateHitCount(path)

		return c.Next()
	}
}

// GinLogger attaches a new contextual logger to the context, for Gin router
func GinLogger(baseLogger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()

		// Create a new logger with the path and IP
		ctxLogger := makeContextLogger(baseLogger, path, c.ClientIP())

		if path != "/favicon.ico" {
			ctxLogger.Info(fmt.Sprintf("req: %s %s", c.Request.Method, c.Request.URL.Path), nil)
		}

		// Set the logger in the context
		c.Set("logger", ctxLogger)

		// Increment the global hit counter
		updateHitCount(path)

		c.Next()
	}
}
