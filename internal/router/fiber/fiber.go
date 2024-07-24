package fiberHandlers

import (
	"fx-service/internal/middleware"
	"fx-service/pkg/config"
	"fx-service/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

type FiberRouter struct {
	Logger *logger.Logger
	Config *config.Config
	App    *fiber.App
}

func (r *FiberRouter) RegisterMiddleware() {
	r.App.Use(middleware.FiberLogger(r.Logger))
	r.App.Use(middleware.FiberRateLimiter(r.Config.RateLimiter))
	r.App.Use(func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		return c.Next()
	})
}

func (r *FiberRouter) RegisterRoutes() {

	// Handle favicon requests from browsers
	r.App.Get("/favicon.ico", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "image/x-icon")
		c.Set("Cache-Control", "public, max-age=86400") // Cache for 1 day (86400 seconds)
		c.Set("Expires", time.Now().Add(24*time.Hour).Format(http.TimeFormat))
		return c.SendFile("./public/favicon.ico")
	})

	// Handle 404
	r.App.All("/", func(c *fiber.Ctx) error {
		return replyError(c, http.StatusNotFound, "Not found")
	})

	// Register the routes
	r.App.Get("/rate/:from/:to", GetRate(r.Config))
	r.App.Get("/rates", GetRates(r.Config)) // ?base=USD&quote=EUR,GBP
	r.App.Get("/status", GetStatus(r.Config))
	r.App.Get("/health", HealthCheck(r.Config))
}

func (r *FiberRouter) Serve(addr string) error {
	return r.App.Listen(addr)
}
